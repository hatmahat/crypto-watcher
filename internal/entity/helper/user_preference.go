package helper

import "crypto-watcher-backend/internal/entity"

type UserAndUserPreference struct {
	entity.User
	PreferenceId int64 `db:"preference_id"`
}
