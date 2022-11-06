package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

//go:generate mockery --name=cacheService --dir . --output ./mocks --exported
type cacheService interface {
	Set(ctx context.Context, userID int64, command enums.CommandType, value string) error
	Get(ctx context.Context, userID int64, command enums.CommandType) (string, error)
	Clear(ctx context.Context, userID int64, command enums.CommandType) error
	ClearKeys(ctx context.Context, userID int64, commands ...enums.CommandType) error
}

type CacheServiceLatencyDecorator struct {
	service cacheService
	latency *prometheus.HistogramVec
}

func NewCacheServiceLatencyDecorator(service cacheService) *CacheServiceLatencyDecorator {
	return &CacheServiceLatencyDecorator{
		service: service,
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "latency_cache_service",
			Help:    "Duration of CacheService methods",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0},
		}, []string{"method"}),
	}
}

func (d *CacheServiceLatencyDecorator) Set(ctx context.Context, userID int64, command enums.CommandType, value string) error {
	startTime := time.Now()
	err := d.service.Set(ctx, userID, command, value)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("Set").Observe(duration.Seconds())

	return err
}

func (d *CacheServiceLatencyDecorator) Get(ctx context.Context, userID int64, command enums.CommandType) (string, error) {
	startTime := time.Now()
	res, err := d.service.Get(ctx, userID, command)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("Set").Observe(duration.Seconds())

	return res, err
}

func (d *CacheServiceLatencyDecorator) Clear(ctx context.Context, userID int64, command enums.CommandType) error {
	startTime := time.Now()
	err := d.service.Clear(ctx, userID, command)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("Clear").Observe(duration.Seconds())

	return err
}

func (d *CacheServiceLatencyDecorator) ClearKeys(ctx context.Context, userID int64, commands ...enums.CommandType) error {
	startTime := time.Now()
	err := d.service.ClearKeys(ctx, userID, commands...)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("ClearKeys").Observe(duration.Seconds())

	return err
}
