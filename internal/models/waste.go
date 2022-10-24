package models

import (
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent"
)

type Waste struct {
	*ent.Waste
}

func NewWaste(category string, cost int64, date time.Time) *Waste {
	return &Waste{
		Waste: &ent.Waste{
			Cost:     cost,
			Category: category,
			Date:     date,
		},
	}
}
