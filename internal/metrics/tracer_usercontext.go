package metrics

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type UserContextServiceTracerDecorator struct {
	service userContextService
	tracer  trace.Tracer
}

func NewUserContextServiceTracerDecorator(service userContextService, tracerProvider *tracesdk.TracerProvider) *UserContextServiceTracerDecorator {
	return &UserContextServiceTracerDecorator{
		service: service,
		tracer:  tracerProvider.Tracer("usercontext-service"),
	}
}

func (d *UserContextServiceTracerDecorator) SetContext(ctx context.Context, userID int64, context enums.UserContext) error {
	ctxTrace, span := d.tracer.Start(ctx, "SetContext")
	defer span.End()

	return d.service.SetContext(ctxTrace, userID, context)
}

func (d *UserContextServiceTracerDecorator) GetContext(ctx context.Context, userID int64) (enums.UserContext, error) {
	ctxTrace, span := d.tracer.Start(ctx, "GetContext")
	defer span.End()

	return d.service.GetContext(ctxTrace, userID)
}

func (d *UserContextServiceTracerDecorator) SetCurrency(ctx context.Context, userID int64, currency string) error {
	ctxTrace, span := d.tracer.Start(ctx, "SetCurrency")
	defer span.End()

	return d.service.SetCurrency(ctxTrace, userID, currency)
}

func (d *UserContextServiceTracerDecorator) GetCurrency(ctx context.Context, userID int64) (string, error) {
	ctxTrace, span := d.tracer.Start(ctx, "GetCurrency")
	defer span.End()

	return d.service.GetCurrency(ctxTrace, userID)
}
