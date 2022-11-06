package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

type WasteRepositoryLatencyDecorator struct {
	wasteRepo wasteRepository
	latency   *prometheus.HistogramVec
}

func NewWasteRepositoryLatencyDecorator(wasteRepo wasteRepository) *WasteRepositoryLatencyDecorator {
	return &WasteRepositoryLatencyDecorator{
		wasteRepo: wasteRepo,
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "latency_waste_repository",
			Help:    "Duration of WasteRepository methods",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0},
		}, []string{"method"}),
	}
}

func (d *WasteRepositoryLatencyDecorator) GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	startTime := time.Now()
	res, err := d.wasteRepo.GetReportLastWeek(ctx, userID)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("GetReportLastWeek").Observe(duration.Seconds())

	return res, err
}

func (d *WasteRepositoryLatencyDecorator) GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	startTime := time.Now()
	res, err := d.wasteRepo.GetReportLastMonth(ctx, userID)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("GetReportLastMonth").Observe(duration.Seconds())

	return res, err
}

func (d *WasteRepositoryLatencyDecorator) GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	startTime := time.Now()
	res, err := d.wasteRepo.GetReportLastYear(ctx, userID)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("GetReportLastYear").Observe(duration.Seconds())

	return res, err
}

func (d *WasteRepositoryLatencyDecorator) SumOfWastesAfterDate(ctx context.Context, userID int64, date time.Time) (int64, error) {
	startTime := time.Now()
	res, err := d.wasteRepo.SumOfWastesAfterDate(ctx, userID, date)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SumOfWastesAfterDate").Observe(duration.Seconds())

	return res, err
}

func (d *WasteRepositoryLatencyDecorator) AddWasteToUser(ctx context.Context, userID int64, waste *models.Waste) (*models.Waste, error) {
	startTime := time.Now()
	res, err := d.wasteRepo.AddWasteToUser(ctx, userID, waste)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("AddWasteToUser").Observe(duration.Seconds())

	return res, err
}
