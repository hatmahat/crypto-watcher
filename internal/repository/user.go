package repository

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/pkg/database"

	"github.com/sirupsen/logrus"
)

type (
	UserRepo interface {
		GetUserByReportTime(ctx context.Context, filter GetUserFilter) ([]entity.User, error)
	}

	userRepo struct {
		db *database.Replication
	}

	UserRepoParam struct {
		DB map[string]*database.Replication
	}

	GetUserFilter struct {
		ReportTime string
		AssetType  string
		AssetCode  string
	}
)

const (
	getUserByReportTime = `select u.id, u.telegram_chat_id from users u
	inner join user_preferences up on up.user_id = u.id
	where up.report_time = $1 and asset_type = $2 and asset_code = $3`
)

func NewUserRepo(param UserRepoParam) UserRepo {
	return &userRepo{
		db: param.DB[config.CryptoWatcherDB],
	}
}

func (ur *userRepo) GetUserByReportTime(ctx context.Context, filter GetUserFilter) ([]entity.User, error) {
	const funcName = "[internal][repository]GetUserByReportTime"
	var (
		result []entity.User
		err    error
	)
	err = ur.db.Select(&result, getUserByReportTime, filter.ReportTime, filter.AssetType, filter.AssetCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err.Error(),
			"report_time": filter.ReportTime,
		}).Errorf("%s: Query Error", funcName)
		return nil, err
	}
	return result, nil
}
