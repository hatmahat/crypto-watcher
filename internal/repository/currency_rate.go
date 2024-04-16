package repository

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/pkg/database"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	CurrencyRateRepo interface {
		GetCurrencyRateByDate(ctx context.Context, currencyPair string, date time.Time) (*entity.CurrencyRate, error)
		InsertCurrencyRate(ctx context.Context, currencyRate entity.CurrencyRate) error
	}

	currencyRateRepo struct {
		db *database.Replication
	}

	CurrencyRateRepoParam struct {
		DB map[string]*database.Replication
	}
)

const (
	getCurrencyRateByDateQuery = `SELECT id, currency_pair, rate, (created_at AT TIME ZONE 'Asia/Jakarta') created_at FROM currency_rates WHERE currency_pair = $1 AND DATE(created_at AT TIME ZONE 'Asia/Jakarta') = $2`
	insertCurrencyRateQuery    = `INSERT INTO currency_rates(currency_pair, rate) VALUES ($1, $2) RETURNING id`
)

func NewCurrencyRateRepo(param CurrencyRateRepoParam) CurrencyRateRepo {
	return &currencyRateRepo{
		db: param.DB[config.CryptoWatcherDB],
	}
}

func (cr *currencyRateRepo) GetCurrencyRateByDate(ctx context.Context, currencyPair string, date time.Time) (*entity.CurrencyRate, error) {
	const funcName = "[internal][repository]GetCurrencyRateByDate"
	var (
		result entity.CurrencyRate
		err    error
	)
	err = cr.db.Get(&result, getCurrencyRateByDateQuery, currencyPair, date.Format("2006-01-02"))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"date": date,
		}).Errorf("%s: Query Error [%s]", funcName, err)
		return nil, err
	}
	return &result, nil
}

func (cr *currencyRateRepo) InsertCurrencyRate(ctx context.Context, currencyRate entity.CurrencyRate) error {
	const funcName = "[internal][repository]InsertCurrencyRate"
	result, err := cr.db.Exec(insertCurrencyRateQuery, currencyRate.CurrencyPair, currencyRate.Rate)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err.Error(),
			"currency_rate": currencyRate,
		}).Errorf("%s: Error Inserting currency_rates [%s]", funcName, err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		return fmt.Errorf("%s: Insert currency_rates failed [%s]", funcName, err)
	}
	return nil
}
