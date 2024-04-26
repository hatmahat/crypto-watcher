package entity

type AssetPrice struct {
	BaseEntity
	AssetType string  `db:"asset_type"`
	AssetCode string  `db:"asset_code"`
	PriceUSD  float64 `db:"price_usd"`
}
