package repository

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/pkg/database"
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	AssetPriceRepo interface {
		InsertAssetPrice(ctx context.Context, assetPrice entity.AssetPrice) error
	}

	assetPriceRepo struct {
		db *database.Replication
	}

	AssetPriceRepoParam struct {
		DB map[string]*database.Replication
	}
)

const (
	insertAssetPriceQuery = `INSERT INTO asset_prices(asset_type, asset_code, price_usd) VALUES ($1, $2, $3)`
)

func NewAssetPriceRepo(param AssetPriceRepoParam) AssetPriceRepo {
	return &assetPriceRepo{
		db: param.DB[config.CryptoWatcherDB],
	}
}

func (ap *assetPriceRepo) InsertAssetPrice(ctx context.Context, assetPrice entity.AssetPrice) error {
	const funcName = "[internal][repository]InsertAssetPrice"
	result, err := ap.db.Exec(insertAssetPriceQuery, assetPrice.AssetType, assetPrice.AssetCode, assetPrice.PriceUSD)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err.Error(),
			"asset_price": assetPrice,
		}).Errorf("%s: Error Inserting asset_price", funcName)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		return fmt.Errorf("%s: Insert currency_rates failed [%s]", funcName, err)
	}
	return nil
}
