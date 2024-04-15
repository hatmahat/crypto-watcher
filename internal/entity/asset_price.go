package entity

type AssetPrice struct {
	BaseEntity
	AssetType string `db:"asset_type"`
	AssetCode string `db:"asset_code"`
	PriceUsd  string `db:"price_usd"`
}
