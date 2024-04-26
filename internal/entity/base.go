package entity

import "time"

type BaseEntity struct {
	Id        int64      `db:"id"`
	Uuid      *string    `db:"uuid"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
