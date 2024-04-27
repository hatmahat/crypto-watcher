package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/whatsapp_cloud_api"
	"net/http"
)

func NewWaMessaging(cfg *config.Config, httpClient *http.Client) whatsapp_cloud_api.WaMessaging {
	return whatsapp_cloud_api.NewWaMessaging(
		cfg.WhatsAppAPIHost,
		cfg.WhatsAppAPIKey,
		cfg.WhatsAppPhoneNumberId,
		httpClient,
	)
}
