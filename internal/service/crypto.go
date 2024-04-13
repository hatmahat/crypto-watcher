package service

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/coin_api"
	"crypto-watcher-backend/pkg/coingecko_api"
	"crypto-watcher-backend/pkg/currency_api"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/whatsapp_cloud_api"
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	CryptoService interface {
		BitcoinPriceWatcher(ctx context.Context) error
	}

	CryptoServiceParam struct {
		CoinGecko   coingecko_api.CoinGecko
		Coin        coin_api.Coin
		Currency    currency_api.Currency
		WaMessaging whatsapp_cloud_api.WaMessaging
		Cfg         *config.Config
	}

	cryptoService struct {
		coinGecko   coingecko_api.CoinGecko
		coin        coin_api.Coin
		currency    currency_api.Currency
		waMessaging whatsapp_cloud_api.WaMessaging
		cfg         *config.Config
	}
)

func NewCryptoService(param CryptoServiceParam) CryptoService {
	return &cryptoService{
		coinGecko:   param.CoinGecko,
		coin:        param.Coin,
		currency:    param.Currency,
		waMessaging: param.WaMessaging,
		cfg:         param.Cfg,
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
			"err": err.Error(),
		}).Errorf("%s: Error Getting Currency Price from Currency API", funcName)
		return err
	}

	usdPrice := format.ThousandSepartor(int64(bitcoinPrice.Bitcoin.USD), ',')
	idrPrice := format.ThousandSepartor(int64(bitcoinPrice.Bitcoin.USD*(*usdToIdr)), '.')
	fmt.Println("USD ", usdPrice)
	fmt.Println("IDR", idrPrice)

	parameters := []string{ // TODO: make parameters dynamic and also add currency conversion to IDR
		"increased",
		"3.5",
		usdPrice,
		idrPrice,
		"up",
		"1,400",
		"yesterday",
		"20.000.000",
		format.GetCurrentTimeInFullFormat()}
	_, err = cs.waMessaging.SendWaMessageByTemplate(ctx, cs.cfg.WhatsAppTestPhoneNumber, whatsapp_cloud_api.BitcoinPriceAlert, parameters)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Sending WA Message", funcName)
		return err
	}

	return nil
}

func (cs cryptoService) convertCurrencyFromUsd(ctx context.Context, currencyCode string) (*int, error) {
	const funcName = "[internal][service]convertCurrencyFromUsd"

	// TODO: need to fix, update currency api only once a day and store it on redis/db

	currencyParams := map[string]string{
		currency_api.Currencies: currencyCode,
	}
	currency, err := cs.currency.GetCurrentCurrency(ctx, currencyParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Currency Price from Currency API", funcName)
		return nil, err
	}

	var usdToIdr int
	if val, ok := currency.Data[currencyCode]; ok {
		usdToIdr = int(val.Value)
	} else {
		return nil, fmt.Errorf("%s: Currency Code [%s] not Found", funcName, currencyCode)
	}

	return &usdToIdr, nil
}
