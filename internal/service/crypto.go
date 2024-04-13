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

	bitcoinParams := map[string]string{
		coingecko_api.Ids:          coingecko_api.Bitcoin,
		coingecko_api.VsCurrencies: coingecko_api.Usd,
	}
	bitcoinPrice, err := cs.coinGecko.GetCurrentPrice(ctx, bitcoinParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Current Price from Coin Gecko", funcName)
		return err
	}

	currencyParams := map[string]string{
		currency_api.Currencies: currency_api.IDR,
	}
	currency, err := cs.currency.GetCurrentCurrency(ctx, currencyParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Currency Price from Currency API", funcName)
		return err
	}

	usdToIdr := int(currency.Data[currency_api.IDR].Value)

	usdPrice := format.ThousandSepartor(int64(bitcoinPrice.Bitcoin.USD), ',')
	idrPrice := format.ThousandSepartor(int64(bitcoinPrice.Bitcoin.USD*usdToIdr), '.')
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
