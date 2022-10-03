package models

import "time"

type Waste struct {
	ID       int64
	Cost     int
	Category string
	Date     time.Time
}

func NewWaste(category string, cost int, date time.Time) *Waste {
	return &Waste{
		Cost:     cost,
		Category: category,
		Date:     date,
	}
}
