package repository

import (
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

type WasteRepositoryInMemory struct {
	wastes map[int64][]models.Waste
	mutex  *sync.RWMutex
}

func NewWasteRepositoryInMemory() *WasteRepositoryInMemory {
	return &WasteRepositoryInMemory{
		wastes: make(map[int64][]models.Waste, 0),
		mutex:  &sync.RWMutex{},
	}
}

func (r *WasteRepositoryInMemory) GetWastesByUser(userID int64) ([]models.Waste, error) {
	r.mutex.RLock()
	wastes, ok := r.wastes[userID]
	r.mutex.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	return wastes, nil
}

func (r *WasteRepositoryInMemory) GetWastesByUserLastWeek(userID int64) ([]models.Waste, error) {
	return r.GetWastesByUserAfterDate(userID, time.Now().Add(-weekDuration))
}

func (r *WasteRepositoryInMemory) GetWastesByUserLastMonth(userID int64) ([]models.Waste, error) {
	return r.GetWastesByUserAfterDate(userID, time.Now().Add(-monthDuration))
}

func (r *WasteRepositoryInMemory) GetWastesByUserLastYear(userID int64) ([]models.Waste, error) {
	return r.GetWastesByUserAfterDate(userID, time.Now().Add(-yearDuration))
}

func (r *WasteRepositoryInMemory) GetWastesByUserAfterDate(userID int64, date time.Time) ([]models.Waste, error) {
	r.mutex.RLock()
	wastes, ok := r.wastes[userID]
	r.mutex.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	result := make([]models.Waste, 0)
	for _, waste := range wastes {
		if waste.Date.After(date) {
			result = append(result, waste)
		}
	}

	if len(result) == 0 {
		return nil, ErrNotFound
	}

	return result, nil
}

func (r *WasteRepositoryInMemory) GetReportLastWeek(userID int64) ([]models.CategoryReport, error) {
	return r.GetReportAfterDate(userID, time.Now().Add(-weekDuration))
}

func (r *WasteRepositoryInMemory) GetReportLastMonth(userID int64) ([]models.CategoryReport, error) {
	return r.GetReportAfterDate(userID, time.Now().Add(-monthDuration))
}

func (r *WasteRepositoryInMemory) GetReportLastYear(userID int64) ([]models.CategoryReport, error) {
	return r.GetReportAfterDate(userID, time.Now().Add(-yearDuration))
}

func (r *WasteRepositoryInMemory) GetReportAfterDate(userID int64, date time.Time) ([]models.CategoryReport, error) {
	r.mutex.RLock()
	wastes, err := r.GetWastesByUserAfterDate(userID, date)
	r.mutex.RUnlock()

	if err != nil {
		return nil, fmt.Errorf("failed to get wastes by user after the date %v: %w", date, err)
	}

	categories := make(map[string]int64)
	for _, waste := range wastes {
		if _, ok := categories[waste.Category]; !ok {
			categories[waste.Category] = 0
		}

		categories[waste.Category] += waste.Cost
	}

	result := make([]models.CategoryReport, 0, len(categories))
	for category, sum := range categories {
		result = append(result, models.CategoryReport{
			Sum:      sum,
			Category: category,
		})
	}

	if len(result) == 0 {
		return nil, ErrNotFound
	}

	return result, nil
}

func (r *WasteRepositoryInMemory) AddWasteToUser(userID int64, waste *models.Waste) error {
	r.mutex.Lock()

	if _, ok := r.wastes[userID]; !ok {
		r.wastes[userID] = make([]models.Waste, 0, 1)
	}

	r.wastes[userID] = append(r.wastes[userID], *waste)

	r.mutex.Unlock()
	return nil
}
