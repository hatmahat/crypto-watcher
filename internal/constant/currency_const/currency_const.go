package currency_const

import "crypto-watcher-backend/pkg/currency_converter_api"

const (
	IDR = "IDR"
	USD = "USD"
)

var CurrencyConverterMapper = map[string]string{
	IDR: currency_converter_api.IDR,
	USD: currency_converter_api.USD,
}

func CurrencyPair(from, to string) string {
	return from + "-" + to
}
