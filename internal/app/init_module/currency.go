package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/currency_api"
	"net/http"
)

func NewCurrency(cfg *config.Config, httpClient *http.Client) currency_api.Currency {
	return currency_api.NewCurrency(
		cfg.CurrencyAPIHost,
		cfg.CurrencyAPIKey,
		httpClient,
	)
}
