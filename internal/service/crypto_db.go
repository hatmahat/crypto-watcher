package service

import (
	"context"
	"crypto-watcher-backend/internal/constant/currency_const"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

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
