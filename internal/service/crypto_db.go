package service

import (
	"context"
	"crypto-watcher-backend/internal/constant/currency_const"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

// convertCurrencyFromUSD converts a given currency from USD to the specified currency code.
// It first attempts to retrieve the currency rate from the database. If the rate is not found,
// it fetches the rate from an external currency converter API and stores it in the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifetime.
//   - currencyCode: The target currency code to which USD should be converted.
//
// Returns:
//   - A pointer to the converted rate as an integer.
//   - An error if any issues occur during the conversion process or data retrieval.
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
