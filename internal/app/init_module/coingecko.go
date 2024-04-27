package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/coingecko_api"
	"net/http"
)

func NewCoinGecko(cfg *config.Config, httpClient *http.Client) coingecko_api.CoinGecko {
	return coingecko_api.NewCoinGecko(
		cfg.CoinGeckoAPIHost,
		httpClient,
	)
}
