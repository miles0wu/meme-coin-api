package domain

import "time"

type Coin struct {
	Id              int64
	Name            string
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PopularityScore uint32
}
