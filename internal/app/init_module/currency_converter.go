package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/currency_converter_api"
	"net/http"
)

func NewCurrencyConverter(cfg *config.Config, httpClient *http.Client) currency_converter_api.CurrencyConverter {
	return currency_converter_api.NewCurrencyConverter(
		cfg.CurrencyConverterAPIHost,
		cfg.CurrencyConverterAPIKey,
		httpClient,
	)
}
