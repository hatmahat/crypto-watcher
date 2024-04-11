package init_module

import (
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/coin_api"
)

func NewCoin(cfg *config.Config) coin_api.Coin {
	return coin_api.NewCoin(
		cfg.CoinAPIHost,
		cfg.CoinAPIKey,
	)
}
