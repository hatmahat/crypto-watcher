package entity

type CurrencyRate struct {
	BaseEntity
	CurrencyPair string  `db:"currency_pair"`
	Rate         float64 `db:"rate"`
}
