package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/currency_converter_api"
)

func NewCurrencyConverter(cfg *config.Config) currency_converter_api.CurrencyConverter {
	return currency_converter_api.NewCurrencyConverter(
		cfg.CurrencyConverterAPIHost,
		cfg.CurrencyConverterAPIKey,
	)
}
