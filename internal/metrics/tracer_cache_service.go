package metrics

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type CacheServiceTracerDecorator struct {
	service cacheService
	tracer  trace.Tracer
}

func NewCacheServiceTracerDecorator(service cacheService, tracerProvider *tracesdk.TracerProvider) *CacheServiceTracerDecorator {
	return &CacheServiceTracerDecorator{
		service: service,
		tracer:  tracerProvider.Tracer("cache-service"),
	}
}

func (d *CacheServiceTracerDecorator) Set(ctx context.Context, userID int64, command enums.CommandType, value string) error {
	ctxTrace, span := d.tracer.Start(ctx, "Set")
	defer span.End()

	return d.service.Set(ctxTrace, userID, command, value)
}

func (d *CacheServiceTracerDecorator) Get(ctx context.Context, userID int64, command enums.CommandType) (string, error) {
	ctxTrace, span := d.tracer.Start(ctx, "Get")
	defer span.End()

	return d.service.Get(ctxTrace, userID, command)
}

func (d *CacheServiceTracerDecorator) Clear(ctx context.Context, userID int64, command enums.CommandType) error {
	ctxTrace, span := d.tracer.Start(ctx, "Clear")
	defer span.End()

	return d.service.Clear(ctxTrace, userID, command)
}

func (d *CacheServiceTracerDecorator) ClearKeys(ctx context.Context, userID int64, commands ...enums.CommandType) error {
	ctxTrace, span := d.tracer.Start(ctx, "ClearKeys")
	defer span.End()

	return d.service.ClearKeys(ctxTrace, userID, commands...)
}
