package repository

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/entity/helper"
	"crypto-watcher-backend/pkg/database"

	"github.com/sirupsen/logrus"
)

type (
	UserRepo interface {
		GetUserAndUserPreferenceByReportTime(ctx context.Context, filter GetUserFilter) ([]helper.UserAndUserPreference, error)
	}

	userRepo struct {
		db *database.Replication
	}

	UserRepoParam struct {
		DB map[string]*database.Replication
	}

	GetUserFilter struct {
		ReportTime     string
		AssetType      string
		AssetCode      string
		PreferenceType string
	}
)

const (
	getUserByReportTimeQuery = `SELECT u.id, u.telegram_chat_id, up.id preference_id FROM users u INNER JOIN user_preferences up ON up.user_id = u.id WHERE up.report_time = $1 AND asset_type = $2 AND asset_code = $3 AND preference_type = $4`
)

func NewUserRepo(param UserRepoParam) UserRepo {
	return &userRepo{
		db: param.DB[config.CryptoWatcherDB],
	}
}

func (ur *userRepo) GetUserAndUserPreferenceByReportTime(ctx context.Context, filter GetUserFilter) ([]helper.UserAndUserPreference, error) {
	const funcName = "[internal][repository]GetUserByReportTime"
	var (
		result []helper.UserAndUserPreference
		err    error
	)
	err = ur.db.Select(&result, getUserByReportTimeQuery, filter.ReportTime, filter.AssetType, filter.AssetCode, filter.PreferenceType)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err.Error(),
			"report_time": filter.ReportTime,
		}).Errorf("%s: Query Error", funcName)
		return nil, err
	}
	return result, nil
}
