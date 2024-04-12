package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/whatsapp_cloud_api"
)

func NewWaMessaging(cfg *config.Config) whatsapp_cloud_api.WaMessaging {
	return whatsapp_cloud_api.NewWaMessaging(
		cfg.WhatsAppAPIHost,
		cfg.WhatsAppAPIKey,
		cfg.WhatsAppPhoneNumberId,
	)
}
