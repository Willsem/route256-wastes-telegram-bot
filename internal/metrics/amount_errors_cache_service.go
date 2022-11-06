package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

type CacheServiceAmountErrorsDecorator struct {
	service     cacheService
	countErrors *prometheus.CounterVec
}

func NewCacheServiceAmountErrorsDecorator(service cacheService) *CacheServiceAmountErrorsDecorator {
	return &CacheServiceAmountErrorsDecorator{
		service: service,
		countErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "count_errors_cache_service",
			Help: "Count of errors in CacheService methods",
		}, []string{"method"}),
	}
}

func (d *CacheServiceAmountErrorsDecorator) Set(ctx context.Context, userID int64, command enums.CommandType, value string) error {
	err := d.service.Set(ctx, userID, command, value)
	if err != nil {
		d.countErrors.WithLabelValues("Set").Inc()
	}
	return err
}

func (d *CacheServiceAmountErrorsDecorator) Get(ctx context.Context, userID int64, command enums.CommandType) (string, error) {
	res, err := d.service.Get(ctx, userID, command)
	if err != nil {
		d.countErrors.WithLabelValues("Get").Inc()
	}
	return res, err
}

func (d *CacheServiceAmountErrorsDecorator) Clear(ctx context.Context, userID int64, command enums.CommandType) error {
	err := d.service.Clear(ctx, userID, command)
	if err != nil {
		d.countErrors.WithLabelValues("Clear").Inc()
	}
	return err
}

func (d *CacheServiceAmountErrorsDecorator) ClearKeys(ctx context.Context, userID int64, commands ...enums.CommandType) error {
	err := d.service.ClearKeys(ctx, userID, commands...)
	if err != nil {
		d.countErrors.WithLabelValues("ClearKeys").Inc()
	}
	return err
}
