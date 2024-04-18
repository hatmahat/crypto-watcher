package currency_converter_api

type (
	CurrencyConverterResponse struct {
		BaseCurrencyCode string          `json:"base_currency_code"`
		BaseCurrencyName string          `json:"base_currency_name"`
		Amount           string          `json:"amount"`
		UpdatedDate      string          `json:"updated_date"`
		Rates            map[string]Rate `json:"rates"`
		Status           string          `json:"status"`
	}

	Rate struct {
		CurrencyName  string `json:"currency_name"`
		Rate          string `json:"rate"`
		RateForAmount string `json:"rate_for_amount"`
	}
)
