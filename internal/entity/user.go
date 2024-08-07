package entity

type User struct {
	BaseEntity
	Username       string  `db:"username"`
	Email          *string `db:"email"`
	PhoneNumber    *string `db:"phone_number"`
	TelegramChatId *int64  `db:"telegram_chat_id"`
}
