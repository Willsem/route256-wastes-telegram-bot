package repository

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/user"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/waste"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

var ErrNotFound = errors.New("wastes not found")

const (
	weekDuration  = time.Hour * 24 * 7
	monthDuration = weekDuration * 4
	yearDuration  = monthDuration * 12
)

type WasteRepository struct {
	client *ent.Client
}

func NewWasteRepository(client *ent.Client) *WasteRepository {
	return &WasteRepository{
		client: client,
	}
}

func (r *WasteRepository) GetWastesByUser(ctx context.Context, userID int64) ([]*models.Waste, error) {
	wastes, err := r.client.Waste.
		Query().
		Where(waste.HasUserWith(user.ID(userID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Waste, 0, len(wastes))
	for _, v := range wastes {
		result = append(result, &models.Waste{
			Waste: v,
		})
	}

	return result, nil
}

func (r *WasteRepository) GetWastesByUserLastWeek(ctx context.Context, userID int64) ([]*models.Waste, error) {
	return r.GetWastesByUserAfterDate(ctx, userID, time.Now().Add(-weekDuration))
}

func (r *WasteRepository) GetWastesByUserLastMonth(ctx context.Context, userID int64) ([]*models.Waste, error) {
	return r.GetWastesByUserAfterDate(ctx, userID, time.Now().Add(-monthDuration))
}

func (r *WasteRepository) GetWastesByUserLastYear(ctx context.Context, userID int64) ([]*models.Waste, error) {
	return r.GetWastesByUserAfterDate(ctx, userID, time.Now().Add(-yearDuration))
}

func (r *WasteRepository) GetWastesByUserAfterDate(
	ctx context.Context, userID int64, date time.Time,
) ([]*models.Waste, error) {
	wastes, err := r.client.Waste.Query().
		Where(waste.HasUserWith(user.ID(userID)), waste.DateGTE(date)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Waste, 0, len(wastes))
	for _, v := range wastes {
		result = append(result, &models.Waste{
			Waste: v,
		})
	}

	return result, nil
}

func (r *WasteRepository) GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	return r.GetReportAfterDate(ctx, userID, time.Now().Add(-weekDuration))
}

func (r *WasteRepository) GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	return r.GetReportAfterDate(ctx, userID, time.Now().Add(-monthDuration))
}

func (r *WasteRepository) GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	return r.GetReportAfterDate(ctx, userID, time.Now().Add(-yearDuration))
}

func (r *WasteRepository) GetReportAfterDate(
	ctx context.Context, userID int64, date time.Time,
) ([]*models.CategoryReport, error) {
	var report []*models.CategoryReport
	err := r.client.Waste.Query().
		Where(waste.HasUserWith(user.ID(userID)), waste.DateGTE(date)).
		GroupBy(waste.FieldCategory).
		Aggregate(ent.Sum(waste.FieldCost)).
		Scan(ctx, &report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *WasteRepository) AddWasteToUser(
	ctx context.Context, userID int64, waste *models.Waste,
) (*models.Waste, error) {
	model, err := r.client.Waste.
		Create().
		SetCost(waste.Cost).
		SetCategory(waste.Category).
		SetDate(waste.Date).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Waste{
		Waste: model,
	}, nil
}

func (r *WasteRepository) SumOfWastesAfterDate(ctx context.Context, userID int64, date time.Time) (int64, error) {
	var result []struct {
		Sum        int64       `json:"sum"`
		UserWastes interface{} `json:"user_wastes"`
	}
	err := r.client.Waste.Query().
		Select(waste.FieldCost).
		Where(waste.HasUserWith(user.ID(userID))).
		GroupBy(waste.UserColumn).
		Aggregate(ent.Sum(waste.FieldCost)).
		Scan(ctx, &result)
	if err != nil {
		return 0, err
	}

	return result[0].Sum, nil
}
