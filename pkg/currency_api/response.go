package currency_api

import "time"

type (
	CurrencyAPIResponse struct {
		Meta Metadata                `json:"meta"`
		Data map[string]CurrencyData `json:"data"`
	}

	Metadata struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	}

	CurrencyData struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	}
)
