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
	"crypto-watcher-backend/pkg/currency_converter_api"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/telegram_bot_api"
	"crypto-watcher-backend/pkg/validation"
	"crypto-watcher-backend/pkg/whatsapp_cloud_api"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	CryptoService interface {
		CryptoWatcher(ctx context.Context) error
	}

	CryptoServiceParam struct {
		CoinGecko         coingecko_api.CoinGecko
		Coin              coin_api.Coin
		CurrencyConverter currency_converter_api.CurrencyConverter
		WaMessaging       whatsapp_cloud_api.WaMessaging
		TelegramBot       telegram_bot_api.TelegramBot
		Cfg               *config.Config
		CurrencyRateRepo  repository.CurrencyRateRepo
		AssetPriceRepo    repository.AssetPriceRepo
	}

	cryptoService struct {
		coinGecko         coingecko_api.CoinGecko
		coin              coin_api.Coin
		currencyConverter currency_converter_api.CurrencyConverter
		waMessaging       whatsapp_cloud_api.WaMessaging
		telegramBot       telegram_bot_api.TelegramBot
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
		telegramBot:       param.TelegramBot,
		cfg:               param.Cfg,
		currencyRateRepo:  param.CurrencyRateRepo,
		assetPriceRepo:    param.AssetPriceRepo,
	}
}

func (cs *cryptoService) CryptoWatcher(ctx context.Context) error {
	const funcName = "[internal][service]CryptoWatcher"

	bitcoinPriceUSD, err := cs.fetchCryptoPriceFromCoinGeckoAPIAndStore(ctx, asset_const.BTC)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Fetching & Storing Bitcoin Price", funcName)
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

	go cs.DailyBitcoinPriceReport(ctx, bitcoinPriceUSD, rateUSDToIDR)

	return nil
}

func (cs *cryptoService) fetchCryptoPriceFromCoinGeckoAPIAndStore(ctx context.Context, assetCode string) (*int, error) {
	const funcName = "[internal][service]fetchCryptoPriceFromCoinGeckoAPIAndStore"

	coinGeckoId, err := validation.ValidateFromMapper(assetCode, asset_const.CoinGeckoMapper)
	if err != nil {
		logrus.Errorf("%s: Asset Code [%s] Not Found", funcName, assetCode)
		return nil, err
	}

	coinGeckoParams := map[string]string{
		coingecko_api.Ids:          *coinGeckoId,
		coingecko_api.VsCurrencies: coingecko_api.USD,
	}
	bitcoinPrice, err := cs.coinGecko.GetCurrentPrice(ctx, coinGeckoParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":               err.Error(),
			"coin_gecko_params": coinGeckoParams,
		}).Errorf("%s: Error Getting Current Price from Coin Gecko", funcName)
		return nil, err
	}

	assetPrice := entity.AssetPrice{
		AssetType: asset_const.CRYPTO,
		AssetCode: assetCode,
		PriceUSD:  float64(bitcoinPrice.Bitcoin.USD),
	}
	err = cs.assetPriceRepo.InsertAssetPrice(ctx, assetPrice)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err.Error(),
			"asset_price": assetPrice,
		}).Errorf("%s: Error Inserting Asset Price", funcName)
		return nil, err
	}

	return &bitcoinPrice.Bitcoin.USD, nil
}

func (cs cryptoService) convertCurrencyFromUSD(ctx context.Context, currencyCode string) (*int, error) {
	const funcName = "[internal][service]convertCurrencyFromUSD"

	currencyPair := currency_const.CurrencyPair(currency_const.USD, currencyCode)
	currencyRate, err := cs.currencyRateRepo.GetCurrencyRateByDate(ctx, currencyPair, time.Now())
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"time": time.Now(),
		}).Errorf("%s: Error Getting Currency Rate from DB", funcName)
		return nil, err
	}

	if currencyRate == nil {
		currencyRate, err = cs.fetchRateFromCurrencyConverterAPIAndStore(ctx, currency_const.USD, currencyCode)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":           err.Error(),
				"currency_code": currencyCode,
			}).Errorf("%s: Error Fetching & Storing Currency Rate", funcName)
			return nil, err
		}
	}

	convertedRate := int(currencyRate.Rate)
	return &convertedRate, nil
}

func (cs *cryptoService) fetchRateFromCurrencyConverterAPIAndStore(ctx context.Context, currencyCodeFrom, currencyCodeTo string) (*entity.CurrencyRate, error) {
	const funcName = "[internal][service]fetchRateFromCurrencyConverterAPIAndStore"

	convertCurrencyFrom, err := validation.ValidateFromMapper(currencyCodeFrom, currency_const.CurrencyConverterMapper)
	if err != nil {
		logrus.Errorf("%s: Currency Converter Code From [%s] Not Found", funcName, currencyCodeTo)
		return nil, err
	}

	convertCurrencyTo, err := validation.ValidateFromMapper(currencyCodeTo, currency_const.CurrencyConverterMapper)
	if err != nil {
		logrus.Errorf("%s: Currency Converter Code To [%s] Not Found", funcName, currencyCodeTo)
		return nil, err
	}

	currencyConverterParams := map[string]string{
		currency_converter_api.Format: currency_converter_api.JSON,
		currency_converter_api.From:   *convertCurrencyFrom,
		currency_converter_api.To:     *convertCurrencyTo,
		currency_converter_api.Amount: "1",
	}
	currencyConverter, err := cs.currencyConverter.GetCurrencyRate(ctx, currencyConverterParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Getting Currency Price from Currency API", funcName)
		return nil, err
	}

	rateStr, ok := currencyConverter.Rates[*convertCurrencyTo]
	if !ok {
		return nil, fmt.Errorf("%s: Currency Code [%s] not Found", funcName, currencyCodeTo)
	}

	rate, err := strconv.ParseFloat(rateStr.Rate, 64)
	if err != nil {
		return nil, fmt.Errorf("%s: Failed to Convert to float64 [%s]", funcName, rateStr.Rate)
	}

	currencyRate := &entity.CurrencyRate{
		Rate:         rate,
		CurrencyPair: currency_const.CurrencyPair(currencyCodeFrom, currencyCodeTo),
	}
	err = cs.currencyRateRepo.InsertCurrencyRate(ctx, *currencyRate)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err.Error(),
			"currency_rate": currencyRate,
		}).Errorf("%s: Error Inserting Currency Rate to DB", funcName)
	}

	return currencyRate, nil
}

func (cs *cryptoService) DailyBitcoinPriceReport(ctx context.Context, bitcoinPriceUSD, rateUSDToIDR *int) {

	// TODO alert every x AM

	// TODO (improvement): not only support bitcoin, get from user preference
	var chatId int64 = 513439237

	// fmt.Println(time.Now().Format("15:04"))

	usdPrice := format.ThousandSepartor(int64(*bitcoinPriceUSD), ',')
	idrPrice := format.ThousandSepartor(int64(*bitcoinPriceUSD*(*rateUSDToIDR)), '.')
	fmt.Printf("USD %s\nIDR %s\n", usdPrice, idrPrice)

	message := telegram_bot_api.BitcoinPriceAlert{
		PercentageIncrease: "3.5",
		CurrentPriceUSD:    usdPrice,
		CurrentPriceIDR:    idrPrice,
		PriceChangeUSD:     "1,400",
		PriceChangeIDR:     "20,000,000",
		FormattedDateTime:  format.GetCurrentTimeInFullFormat(),
	}
	go cs.sendTelegramMessage(chatId, &message)
}

func (cs *cryptoService) sendTelegramMessage(chatId int64, message telegram_bot_api.Message) {
	const funcName = "[internal][service]sendTelegramMessage"
	err := cs.telegramBot.SendTelegramMessageByMessageId(chatId, message.Message())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"message": message.Message(),
			"chat_id": chatId,
		}).Errorf("%s: Error Sending Message Via Telegram", funcName)
	}
}
