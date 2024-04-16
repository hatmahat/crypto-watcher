package entity

import "time"

type CurrencyRate struct {
	Id           int64     `db:"id"`
	CurrencyPair string    `db:"currency_pair"`
	Rate         float64   `db:"rate"`
	CreatedAt    time.Time `db:"created_at"`
}
