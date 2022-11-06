package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

type UserContextServiceLatencyDecorator struct {
	service userContextService
	latency *prometheus.HistogramVec
}

func NewUserContextServiceLatencyDecorator(service userContextService) *UserContextServiceLatencyDecorator {
	return &UserContextServiceLatencyDecorator{
		service: service,
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "latency_usercontext_service",
			Help:    "Duration of UserContextService methods",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0},
		}, []string{"method"}),
	}
}

func (d *UserContextServiceLatencyDecorator) SetContext(ctx context.Context, userID int64, context enums.UserContext) error {
	startTime := time.Now()
	err := d.service.SetContext(ctx, userID, context)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SetContext").Observe(duration.Seconds())

	return err
}

func (d *UserContextServiceLatencyDecorator) GetContext(ctx context.Context, userID int64) (enums.UserContext, error) {
	startTime := time.Now()
	res, err := d.service.GetContext(ctx, userID)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("GetContext").Observe(duration.Seconds())

	return res, err
}

func (d *UserContextServiceLatencyDecorator) SetCurrency(ctx context.Context, userID int64, currency string) error {
	startTime := time.Now()
	err := d.service.SetCurrency(ctx, userID, currency)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SetCurrency").Observe(duration.Seconds())

	return err
}

func (d *UserContextServiceLatencyDecorator) GetCurrency(ctx context.Context, userID int64) (string, error) {
	startTime := time.Now()
	res, err := d.service.GetCurrency(ctx, userID)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("GetCurrency").Observe(duration.Seconds())

	return res, err
}
