package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/telegram_bot_api"
)

func NewTelegramBot(cfg *config.Config) telegram_bot_api.TelegramBot {
	return telegram_bot_api.NewTelegramBot(
		cfg.TelegramBotAPIKey,
	)
}
