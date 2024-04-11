package coin_api

import "time"

type CoinRateResponse struct {
	Time         time.Time `json:"time"`
	AssetIDBase  string    `json:"asset_id_base"`
	AssetIDQuote string    `json:"asset_id_quote"`
	Rate         float64   `json:"rate"`
}
