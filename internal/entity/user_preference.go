package entity

type UserPreference struct {
	BaseEntity
	UserId              int64   `db:"user_id"`
	AssetType           string  `db:"asset_type"`
	AssetCode           string  `db:"asset_code"`
	ThresholdPercentage float64 `db:"threshold_percentage"`
	ObservationPeriod   int64   `db:"observation_period"`
}
