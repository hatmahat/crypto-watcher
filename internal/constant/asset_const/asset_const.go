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

	// Coin name
	Bitcoin  = "Bitcoin"
	Ethereum = "Ethereum"
)

var CoinGeckoMapper = map[string]string{
	BTC: coingecko_api.Bitcoin,
	ETH: coingecko_api.Ethereum,
}

var Coins = []string{
	BTC,
	ETH,
}

var CoinNameMapper = map[string]string{
	BTC: Bitcoin,
	ETH: Ethereum,
}
