package asset_const

import "crypto-watcher-backend/pkg/coingecko_api"

const (
	// Asset Type
	CRYPTO = "CRYPTO"
	STOCK  = "STOCK"

	// Crypto Asset Code
	BTC = "BTC"
	ETH = "ETH"

	// Stock Asset Code
	AAPL  = "AAPL"
	GOOGL = "GOOGL"
)

var CoinGeckoMapper = map[string]string{
	BTC: coingecko_api.Bitcoin,
}
