package repository

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/database"

	"github.com/sirupsen/logrus"
)

type (
	UserPreferenceRepo interface {
		GetDistinctUserPreferenceAssetCodeByAssetType(ctx context.Context, assetType string) ([]string, error)
	}

	userPreferenceRepo struct {
		db *database.Replication
	}

	UserPreferenceRepoParam struct {
		DB map[string]*database.Replication
	}
)

const (
	getDistinctUserPreferenceAssetCodeByAssetTypeQuery = `SELECT DISTINCT(asset_code) FROM user_preferences up WHERE asset_type = $1`
)

func NewUserPreferenceRepo(param UserPreferenceRepoParam) UserPreferenceRepo {
	return &userPreferenceRepo{
		db: param.DB[config.CryptoWatcherDB],
	}
}

func (ur *userPreferenceRepo) GetDistinctUserPreferenceAssetCodeByAssetType(ctx context.Context, assetType string) ([]string, error) {
	const funcName = "[internal][repository]GetDistinctUserPreferenceAssetCodeByAssetType"
	var (
		result []string
		err    error
	)
	err = ur.db.Select(&result, getDistinctUserPreferenceAssetCodeByAssetTypeQuery, assetType)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err.Error(),
			"asset_type": assetType,
		}).Errorf("%s: Query Error", funcName)
		return nil, err
	}
	return result, nil
}
