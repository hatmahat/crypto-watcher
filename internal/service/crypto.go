package service

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/constant/asset_const"
	"crypto-watcher-backend/internal/constant/currency_const"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/pkg/coingecko_api"
	"crypto-watcher-backend/pkg/currency_converter_api"
	"crypto-watcher-backend/pkg/telegram_bot_api"

	"github.com/sirupsen/logrus"
)

type (
	CryptoService interface {
		CryptoWatcher(ctx context.Context) error
	}

	CryptoServiceParam struct {
		Cfg                *config.Config
		CoinGecko          coingecko_api.CoinGecko
		CurrencyConverter  currency_converter_api.CurrencyConverter
		TelegramBot        telegram_bot_api.TelegramBot
		CurrencyRateRepo   repository.CurrencyRateRepo
		AssetPriceRepo     repository.AssetPriceRepo
		UserRepo           repository.UserRepo
		NotifRepo          repository.NotificationRepo
		UserPreferenceRepo repository.UserPreferenceRepo
	}

	cryptoService struct {
		cfg                *config.Config
		coinGecko          coingecko_api.CoinGecko
		currencyConverter  currency_converter_api.CurrencyConverter
		telegramBot        telegram_bot_api.TelegramBot
		currencyRateRepo   repository.CurrencyRateRepo
		assetPriceRepo     repository.AssetPriceRepo
		userRepo           repository.UserRepo
		notifRepo          repository.NotificationRepo
		userPreferenceRepo repository.UserPreferenceRepo
	}
)

func NewCryptoService(param CryptoServiceParam) CryptoService {
	return &cryptoService{
		cfg:                param.Cfg,
		coinGecko:          param.CoinGecko,
		currencyConverter:  param.CurrencyConverter,
		telegramBot:        param.TelegramBot,
		currencyRateRepo:   param.CurrencyRateRepo,
		assetPriceRepo:     param.AssetPriceRepo,
		userRepo:           param.UserRepo,
		notifRepo:          param.NotifRepo,
		userPreferenceRepo: param.UserPreferenceRepo,
	}
}

func (cs *cryptoService) CryptoWatcher(ctx context.Context) error {
	const funcName = "[internal][service]CryptoWatcher"

	assetCodes, err := cs.userPreferenceRepo.GetDistinctUserPreferenceAssetCodeByAssetType(ctx, asset_const.CRYPTO)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Distinct Coin from user_preference", funcName)
		return err
	}

	rateUSDToIDR, err := cs.convertCurrencyFromUSD(ctx, currency_const.IDR)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err.Error(),
			"currency_code": currency_const.IDR,
		}).Errorf("%s: Error Getting Currency Price from Currency API", funcName)
		return err
	}

	assetPrices, err := cs.fetchCryptoPriceFromCoinGeckoAPIAndStore(ctx, assetCodes)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Fetching & Storing Bitcoin Price", funcName)
		return err
	}

	for _, assetPrice := range assetPrices {
		go cs.dailyCoinPriceReport(ctx, assetPrice.AssetCode, assetPrice.PriceUSD, *rateUSDToIDR)
	}

	return nil
}
