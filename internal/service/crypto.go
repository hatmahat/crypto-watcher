package service

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/constant/currency_const"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/pkg/coin_api"
	"crypto-watcher-backend/pkg/coingecko_api"
	"crypto-watcher-backend/pkg/currency_api"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/whatsapp_cloud_api"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	CryptoService interface {
		BitcoinPriceWatcher(ctx context.Context) error
	}

	CryptoServiceParam struct {
		CoinGecko        coingecko_api.CoinGecko
		Coin             coin_api.Coin
		Currency         currency_api.Currency
		WaMessaging      whatsapp_cloud_api.WaMessaging
		Cfg              *config.Config
		CurrencyRateRepo repository.CurrencyRateRepo
	}

	cryptoService struct {
		coinGecko        coingecko_api.CoinGecko
		coin             coin_api.Coin
		currency         currency_api.Currency
		waMessaging      whatsapp_cloud_api.WaMessaging
		cfg              *config.Config
		currencyRateRepo repository.CurrencyRateRepo
	}
)

func NewCryptoService(param CryptoServiceParam) CryptoService {
	return &cryptoService{
		coinGecko:        param.CoinGecko,
		coin:             param.Coin,
		currency:         param.Currency,
		waMessaging:      param.WaMessaging,
		cfg:              param.Cfg,
		currencyRateRepo: param.CurrencyRateRepo,
	}
}

func (cs *cryptoService) BitcoinPriceWatcher(ctx context.Context) error {
	const funcName = "[internal][service]BitcoinPriceWatcher"

	coinGeckoParams := map[string]string{
		coingecko_api.Ids:          coingecko_api.Bitcoin,
		coingecko_api.VsCurrencies: coingecko_api.Usd,
	}
	bitcoinPrice, err := cs.coinGecko.GetCurrentPrice(ctx, coinGeckoParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Current Price from Coin Gecko", funcName)
		return err
	}

	usdToIdr, err := cs.convertCurrencyFromUsd(ctx, currency_api.IDR)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err.Error(),
			"currency_code": currency_api.IDR,
		}).Errorf("%s: Error Getting Currency Price from Currency API", funcName)
		return err
	}

	usdPrice := format.ThousandSepartor(int64(bitcoinPrice.Bitcoin.USD), ',')
	idrPrice := format.ThousandSepartor(int64(bitcoinPrice.Bitcoin.USD*(*usdToIdr)), '.')
	fmt.Printf("USD %s\nIDR %s\n", usdPrice, idrPrice)

	// parameters := []string{ // TODO: make parameters dynamic
	// 	"increased",
	// 	"3.5",
	// 	usdPrice,
	// 	idrPrice,
	// 	"up",
	// 	"1,400",
	// 	"yesterday",
	// 	"20.000.000",
	// 	format.GetCurrentTimeInFullFormat()}
	// _, err = cs.waMessaging.SendWaMessageByTemplate(ctx, cs.cfg.WhatsAppTestPhoneNumber, whatsapp_cloud_api.BitcoinPriceAlert, parameters)
	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"err": err.Error(),
	// 	}).Errorf("%s: Error Sending WA Message", funcName)
	// 	return err
	// }

	return nil
}

func (cs cryptoService) convertCurrencyFromUsd(ctx context.Context, currencyCode string) (*int, error) {
	const funcName = "[internal][service]convertCurrencyFromUsd"

	currencyPair, err := validateCurrencyCode(currencyCode)
	if err != nil {
		logrus.Errorf("%s: Currency Pair [%s] Not Found", funcName, currencyCode)
		return nil, fmt.Errorf("%s: Currency Pair [%s] Not Found", funcName, currencyCode)
	}

	currencyRate, err := cs.currencyRateRepo.GetCurrencyRateByDate(ctx, *currencyPair, time.Now())
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"time": time.Now(),
		}).Errorf("%s: Error Getting Currency Rate from DB [%s]", funcName, err)
		return nil, err
	}

	if currencyRate == nil {
		currencyRate, err = cs.fetchRateFromCurrencyAPIAndStore(ctx, currencyCode)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":           err.Error(),
				"currency_code": currencyCode,
			}).Errorf("%s: Error Fetching & Storing Currency Rate [%s]", funcName, err)
			return nil, err
		}
	}

	convertedRate := int(currencyRate.Rate)
	return &convertedRate, nil
}

func (cs *cryptoService) fetchRateFromCurrencyAPIAndStore(ctx context.Context, currencyCode string) (*entity.CurrencyRate, error) {
	const funcName = "[internal][service]fetchRateFromCurrencyAPIAndStore"

	currencyPair, err := validateCurrencyCode(currencyCode)
	if err != nil {
		logrus.Errorf("%s: Currency Pair [%s] Not Found", funcName, currencyCode)
		return nil, fmt.Errorf("%s: Currency Pair [%s] Not Found", funcName, currencyCode)
	}

	currencyParams := map[string]string{
		currency_api.Currencies: currencyCode,
	}
	currency, err := cs.currency.GetCurrentCurrency(ctx, currencyParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Currency Price from Currency API [%s]", funcName, err)
		return nil, err
	}
	val, ok := currency.Data[currencyCode]
	if !ok {
		return nil, fmt.Errorf("%s: Currency Code [%s] not Found", funcName, currencyCode)
	}
	currencyRate := &entity.CurrencyRate{
		Rate:         val.Value,
		CurrencyPair: *currencyPair,
	}
	err = cs.currencyRateRepo.InsertCurrencyRate(ctx, *currencyRate)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Inserting Currency Rate to DB [%s]", funcName, err)
	}

	return currencyRate, nil
}

func validateCurrencyCode(currencyCode string) (*string, error) {
	const funcName = "[internal][service]validateCurrencyCodes"

	currencyPair, ok := currency_const.CurrencyPairMapper[currencyCode]
	if !ok {
		logrus.Errorf("%s: Currency Pair [%s] Not Found", funcName, currencyCode)
		return nil, fmt.Errorf("%s: Currency Pair [%s] Not Found", funcName, currencyCode)
	}
	return &currencyPair, nil
}
