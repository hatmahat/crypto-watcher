package entity

type Notification struct {
	BaseEntity
	UserId       int64  `db:"user_id"`
	PreferenceId int64  `db:"preference_id"`
	Parameters   []byte `db:"parameters"`
}
