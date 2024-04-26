package entity

import "time"

type UserPreference struct {
	BaseEntity
	UserId              int64      `db:"user_id"`
	PreferenceType      string     `db:"preference_type"`
	AssetType           string     `db:"asset_type"`
	AssetCode           string     `db:"asset_code"`
	ThresholdPercentage *float64   `db:"threshold_percentage"`
	ObservationPeriod   int64      `db:"observation_period"`
	ReportTime          *time.Time `db:"report_time"`
	IsActive            bool       `db:"is_active"`
}
