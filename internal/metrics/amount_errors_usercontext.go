package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

//go:generate mockery --name=userContextService --dir . --output ./mocks --exported
type userContextService interface {
	SetContext(ctx context.Context, userID int64, context enums.UserContext) error
	GetContext(ctx context.Context, userID int64) (enums.UserContext, error)
	SetCurrency(ctx context.Context, userID int64, currency string) error
	GetCurrency(ctx context.Context, userID int64) (string, error)
}

type UserContextServiceAmountErrorsDecorator struct {
	service     userContextService
	countErrors *prometheus.CounterVec
}

func NewUserContextServiceAmountErrorsDecorator(service userContextService) *UserContextServiceAmountErrorsDecorator {
	return &UserContextServiceAmountErrorsDecorator{
		service: service,
		countErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "count_errors_usercontext_service",
			Help: "Count of errors in UserContextService methods",
		}, []string{"method"}),
	}
}

func (d *UserContextServiceAmountErrorsDecorator) SetContext(ctx context.Context, userID int64, context enums.UserContext) error {
	err := d.service.SetContext(ctx, userID, context)
	if err != nil {
		d.countErrors.WithLabelValues("SetContext").Inc()
	}
	return err
}

func (d *UserContextServiceAmountErrorsDecorator) GetContext(ctx context.Context, userID int64) (enums.UserContext, error) {
	res, err := d.service.GetContext(ctx, userID)
	if err != nil {
		d.countErrors.WithLabelValues("GetContext").Inc()
	}
	return res, err
}

func (d *UserContextServiceAmountErrorsDecorator) SetCurrency(ctx context.Context, userID int64, currency string) error {
	err := d.service.SetCurrency(ctx, userID, currency)
	if err != nil {
		d.countErrors.WithLabelValues("SetCurrency").Inc()
	}
	return err
}

func (d *UserContextServiceAmountErrorsDecorator) GetCurrency(ctx context.Context, userID int64) (string, error) {
	res, err := d.service.GetCurrency(ctx, userID)
	if err != nil {
		d.countErrors.WithLabelValues("GetCurrency").Inc()
	}
	return res, err
}
