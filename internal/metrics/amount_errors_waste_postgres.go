package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

//go:generate mockery --name=wasteRepository --dir . --output ./mocks --exported
type wasteRepository interface {
	GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	SumOfWastesAfterDate(ctx context.Context, userID int64, date time.Time) (int64, error)

	AddWasteToUser(ctx context.Context, userID int64, waste *models.Waste) (*models.Waste, error)
}

type WasteRepositoryAmountErrorsDecorator struct {
	wasteRepo   wasteRepository
	countErrors *prometheus.CounterVec
}

func NewWasteRepositoryAmountErrorsDecorator(wasteRepo wasteRepository) *WasteRepositoryAmountErrorsDecorator {
	return &WasteRepositoryAmountErrorsDecorator{
		wasteRepo: wasteRepo,
		countErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "count_errors_waste_repository",
			Help: "Count of errors in WasteRepository methods",
		}, []string{"method"}),
	}
}

func (d *WasteRepositoryAmountErrorsDecorator) GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	res, err := d.wasteRepo.GetReportLastWeek(ctx, userID)
	if err != nil {
		d.countErrors.WithLabelValues("GetReportLastWeek").Inc()
	}
	return res, err
}

func (d *WasteRepositoryAmountErrorsDecorator) GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	res, err := d.wasteRepo.GetReportLastMonth(ctx, userID)
	if err != nil {
		d.countErrors.WithLabelValues("GetReportLastMonth").Inc()
	}
	return res, err
}

func (d *WasteRepositoryAmountErrorsDecorator) GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	res, err := d.wasteRepo.GetReportLastYear(ctx, userID)
	if err != nil {
		d.countErrors.WithLabelValues("GetReportLastYear").Inc()
	}
	return res, err
}

func (d *WasteRepositoryAmountErrorsDecorator) SumOfWastesAfterDate(ctx context.Context, userID int64, date time.Time) (int64, error) {
	res, err := d.wasteRepo.SumOfWastesAfterDate(ctx, userID, date)
	if err != nil {
		d.countErrors.WithLabelValues("SumOfWastesAfterDate").Inc()
	}
	return res, err
}

func (d *WasteRepositoryAmountErrorsDecorator) AddWasteToUser(ctx context.Context, userID int64, waste *models.Waste) (*models.Waste, error) {
	res, err := d.wasteRepo.AddWasteToUser(ctx, userID, waste)
	if err != nil {
		d.countErrors.WithLabelValues("AddWasteToUser").Inc()
	}
	return res, err
}
