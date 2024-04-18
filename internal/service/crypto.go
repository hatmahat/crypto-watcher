package service

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/constant/asset_const"
	"crypto-watcher-backend/internal/constant/currency_const"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/pkg/coin_api"
	"crypto-watcher-backend/pkg/coingecko_api"
	"crypto-watcher-backend/pkg/currency_api"
	"crypto-watcher-backend/pkg/currency_converter_api"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/whatsapp_cloud_api"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	CryptoService interface {
		BitcoinWatcher(ctx context.Context) error
	}

	CryptoServiceParam struct {
		CoinGecko         coingecko_api.CoinGecko
		Coin              coin_api.Coin
		CurrencyConverter currency_converter_api.CurrencyConverter
		WaMessaging       whatsapp_cloud_api.WaMessaging
		Cfg               *config.Config
		CurrencyRateRepo  repository.CurrencyRateRepo
		AssetPriceRepo    repository.AssetPriceRepo
	}

	cryptoService struct {
		coinGecko         coingecko_api.CoinGecko
		coin              coin_api.Coin
		currencyConverter currency_converter_api.CurrencyConverter
		waMessaging       whatsapp_cloud_api.WaMessaging
		cfg               *config.Config
		currencyRateRepo  repository.CurrencyRateRepo
		assetPriceRepo    repository.AssetPriceRepo
	}
)

func NewCryptoService(param CryptoServiceParam) CryptoService {
	return &cryptoService{
		coinGecko:         param.CoinGecko,
		coin:              param.Coin,
		currencyConverter: param.CurrencyConverter,
		waMessaging:       param.WaMessaging,
		cfg:               param.Cfg,
		currencyRateRepo:  param.CurrencyRateRepo,
		assetPriceRepo:    param.AssetPriceRepo,
	}
}

func (cs *cryptoService) BitcoinWatcher(ctx context.Context) error {
	const funcName = "[internal][service]BitcoinWatcher"

	bitcoinPrice, err := cs.fetchBitcoinPriceFromCoinGeckoAPIAndStore(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Fetching & Storing Bitcoin Price [%s]", funcName, err)
		return err
	}

	usdToIdr, err := cs.convertCurrencyFromUsd(ctx, currency_api.IDR)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err.Error(),
			"currency_code": currency_api.IDR,
		}).Errorf("%s: Error Getting Currency Price from Currency API [%s]", funcName, err)
		return err
	}

	usdPrice := format.ThousandSepartor(int64(*bitcoinPrice), ',')
	idrPrice := format.ThousandSepartor(int64(*bitcoinPrice*(*usdToIdr)), '.')
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

func (cs *cryptoService) fetchBitcoinPriceFromCoinGeckoAPIAndStore(ctx context.Context) (*int, error) {
	const funcName = "[internal][service]fetchBitcoinPriceFromCoinGeckoAPIAndStore"

	// TODO (improvement): not only support bitcoin, get from user preference
	coinGeckoParams := map[string]string{
		coingecko_api.Ids:          coingecko_api.Bitcoin,
		coingecko_api.VsCurrencies: coingecko_api.Usd,
	}
	bitcoinPrice, err := cs.coinGecko.GetCurrentPrice(ctx, coinGeckoParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Current Price from Coin Gecko [%s]", funcName, err)
		return nil, err
	}

	assetPrice := entity.AssetPrice{
		AssetType: asset_const.CRYPTO,
		AssetCode: asset_const.BTC,
		PriceUsd:  float64(bitcoinPrice.Bitcoin.USD),
	}

	err = cs.assetPriceRepo.InsertAssetPrice(ctx, assetPrice)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err.Error(),
			"asset_price": assetPrice,
		}).Errorf("%s: Error Inserting Asset Price [%s]", funcName, err)
		return nil, err
	}

	return &bitcoinPrice.Bitcoin.USD, nil
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

	currencyConverterParams := map[string]string{
		currency_converter_api.Format: currency_converter_api.JSON,
		currency_converter_api.From:   currency_converter_api.USD,
		currency_converter_api.To:     currency_converter_api.IDR,
		currency_converter_api.Amount: "1",
	}
	currencyConverter, err := cs.currencyConverter.GetCurrencyRate(ctx, currencyConverterParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Currency Price from Currency API [%s]", funcName, err)
		return nil, err
	}

	rateStr, ok := currencyConverter.Rates[currency_converter_api.IDR]
	if !ok {
		return nil, fmt.Errorf("%s: Currency Code [%s] not Found", funcName, currencyCode)
	}

	rate, err := strconv.ParseFloat(rateStr.Rate, 64)
	if err != nil {
		return nil, fmt.Errorf("%s: Failed to Convert to float64 [%s]", funcName, rateStr.Rate)
	}

	currencyRate := &entity.CurrencyRate{
		Rate:         rate,
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
