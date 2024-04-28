package entity

type Notification struct {
	BaseEntity
	UserId       int64     `db:"user_id"`
	PreferenceId int64     `db:"preference_id"`
	Status       string    `db:"status"`
	Metadata     *Metadata `db:"metadata"`
}

type Metadata struct {
	Message  *string `json:"message,omitempty"`
	Provider *string `json:"provider,omitempty"`
}
