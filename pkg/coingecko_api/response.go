package coingecko_api

type CoinGeckoPriceResponse struct {
	Bitcoin Currency `json:"bitcoin"`
}

type Currency struct {
	USD int `json:"usd"`
}
