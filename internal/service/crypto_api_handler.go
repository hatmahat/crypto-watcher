package service

import (
	"context"
	"crypto-watcher-backend/internal/constant/asset_const"
	"crypto-watcher-backend/internal/constant/currency_const"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/pkg/coingecko_api"
	"crypto-watcher-backend/pkg/currency_converter_api"
	"crypto-watcher-backend/pkg/validation"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

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
