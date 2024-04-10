package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/whatsapp"
)

func NewWaMessaging(cfg *config.Config) whatsapp.WaMessaging {
	return whatsapp.NewWaMessaging(
		cfg.WhatsAppHost,
		cfg.WhatsAppAPIKey,
		cfg.WhatsAppPhoneNumberId,
	)
}
