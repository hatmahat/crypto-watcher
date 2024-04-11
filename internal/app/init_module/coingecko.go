package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/coingecko_api"
)

func NewCoinGecko(cfg *config.Config) coingecko_api.CoinGecko {
	return coingecko_api.NewCoinGecko(
		cfg.CoinGeckoAPIHost,
	)
}
