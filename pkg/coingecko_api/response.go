package coingecko_api

type (
	CoinGeckoPriceResponse map[string]Currency

	Currency struct {
		USD float64 `json:"usd"`
	}
)
