package entity

import "time"

type AssetPrice struct {
	Id        int64     `db:"id"`
	AssetType string    `db:"asset_type"`
	AssetCode string    `db:"asset_code"`
	PriceUsd  float64   `db:"price_usd"`
	CreatedAt time.Time `db:"created_at"`
}
