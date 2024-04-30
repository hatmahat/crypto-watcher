package coingecko_api

type (
	CoinGeckoPriceResponse struct {
		Bitcoin  Currency `json:"bitcoin"`
		Ethereum Currency `json:"ethereum"` // TODO (improvement)
	}

	Currency struct {
		USD int `json:"usd"`
	}
)
