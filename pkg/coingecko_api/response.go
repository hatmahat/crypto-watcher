package coingecko_api

type (
	CoinGeckoPriceResponse struct {
		Bitcoin Currency `json:"bitcoin"`
	}

	Currency struct {
		USD int `json:"usd"`
	}
)
