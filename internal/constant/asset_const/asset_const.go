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

var CoinGeckoMapperToAssetCode = map[string]string{
	coingecko_api.Bitcoin:  BTC,
	coingecko_api.Ethereum: ETH,
}

var AssetCodes = []string{
	BTC,
	ETH,
}

var AssetCodeNameMapper = map[string]string{
	BTC: Bitcoin,
	ETH: Ethereum,
}
