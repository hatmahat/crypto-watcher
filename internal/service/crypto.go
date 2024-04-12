package service

import (
	"context"
	"crypto-watcher-backend/pkg/coin_api"
	"crypto-watcher-backend/pkg/coingecko_api"
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
		WaMessaging whatsapp_cloud_api.WaMessaging
	}

	cryptoService struct {
		coinGecko   coingecko_api.CoinGecko
		coin        coin_api.Coin
		waMessaging whatsapp_cloud_api.WaMessaging
	}
)

func NewCryptoService(param CryptoServiceParam) CryptoService {
	return &cryptoService{
		coinGecko:   param.CoinGecko,
		coin:        param.Coin,
		waMessaging: param.WaMessaging,
	}
}

func (cs *cryptoService) BitcoinPriceWatcher(ctx context.Context) error {
	const funcName = "[internal][service]BitcoinPriceWatcher"

	params := map[string]string{
		coingecko_api.Ids:          coingecko_api.Bitcoin,
		coingecko_api.VsCurrencies: coingecko_api.Usd,
	}
	bitcoinPrice, err := cs.coinGecko.GetCurrentPrice(ctx, params)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("Error Getting Current Price from Coin Gecko: %s", funcName)
		return err
	}

	fmt.Println("PRICE", bitcoinPrice.Bitcoin.USD)

	return nil
}
