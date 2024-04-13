package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/currency_api"
)

func NewCurrency(cfg *config.Config) currency_api.Currency {
	return currency_api.NewCurrency(
		cfg.CurrencyAPIHost,
		cfg.CurrencyAPIKey,
	)
}
