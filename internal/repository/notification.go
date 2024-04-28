package repository

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/pkg/database"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	NotificationRepo interface {
		InsertNotification(ctx context.Context, notification entity.Notification) error
	}

	notificationRepo struct {
		db *database.Replication
	}

	NotificationRepoParam struct {
		DB map[string]*database.Replication
	}
)

const (
	insertNotificationQuery = `INSERT INTO notifications(user_id, preference_id, status, metadata) VALUES ($1, $2, $3, $4)`
)

func NewNotificationRepo(param NotificationRepoParam) NotificationRepo {
	return &notificationRepo{
		db: param.DB[config.CryptoWatcherDB],
	}
}

func (nr *notificationRepo) InsertNotification(ctx context.Context, notification entity.Notification) error {
	const funcName = "[internal][repository]InsertNotification"

	parametersJSON, err := json.Marshal(notification.Metadata)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":          err.Error(),
			"notification": notification,
		}).Errorf("%s: Error marshalling parameters to JSON", funcName)
		return err
	}

	result, err := nr.db.Exec(insertNotificationQuery, notification.UserId, notification.PreferenceId, notification.Status, parametersJSON)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":          err.Error(),
			"notification": notification,
		}).Errorf("%s: Error Inserting notifications", funcName)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		return fmt.Errorf("%s: Insert notifications failed [%d]", funcName, rowsAffected)
	}
	return nil
}
