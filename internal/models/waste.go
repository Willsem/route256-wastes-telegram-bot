package models

import "time"

type Waste struct {
	ID       int64
	Cost     int64
	Category string
	Date     time.Time
}

func NewWaste(category string, cost int64, date time.Time) *Waste {
	return &Waste{
		Cost:     cost,
		Category: category,
		Date:     date,
	}
}
