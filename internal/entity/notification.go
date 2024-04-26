package entity

type Notification struct {
	BaseEntity
	UserId       int64  `db:"user_id"`
	PreferenceId int64  `db:"preference_id"`
	Status       string `db:"status"`
	Parameters   []byte `db:"parameters"`
}
