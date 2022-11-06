package metrics

import (
	"context"
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type WasteRepositoryTracerDecorator struct {
	wasteRepo wasteRepository
	tracer    trace.Tracer
}

func NewWasteRepositoryTracerDecorator(wasteRepo wasteRepository, tracerProvider *tracesdk.TracerProvider) *WasteRepositoryTracerDecorator {
	return &WasteRepositoryTracerDecorator{
		wasteRepo: wasteRepo,
		tracer:    tracerProvider.Tracer("waste-repository"),
	}
}

func (d *WasteRepositoryTracerDecorator) GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	ctxTrace, span := d.tracer.Start(ctx, "GetReportLstWeek")
	defer span.End()

	return d.wasteRepo.GetReportLastWeek(ctxTrace, userID)
}

func (d *WasteRepositoryTracerDecorator) GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	ctxTrace, span := d.tracer.Start(ctx, "GetReportLastMonth")
	defer span.End()

	return d.wasteRepo.GetReportLastMonth(ctxTrace, userID)
}

func (d *WasteRepositoryTracerDecorator) GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error) {
	ctxTrace, span := d.tracer.Start(ctx, "GetReportLastYear")
	defer span.End()

	return d.wasteRepo.GetReportLastYear(ctxTrace, userID)
}

func (d *WasteRepositoryTracerDecorator) SumOfWastesAfterDate(ctx context.Context, userID int64, date time.Time) (int64, error) {
	ctxTrace, span := d.tracer.Start(ctx, "SumOfWastesAfterDate")
	defer span.End()

	return d.wasteRepo.SumOfWastesAfterDate(ctxTrace, userID, date)
}

func (d *WasteRepositoryTracerDecorator) AddWasteToUser(ctx context.Context, userID int64, waste *models.Waste) (*models.Waste, error) {
	ctxTrace, span := d.tracer.Start(ctx, "AddWasteToUser")
	defer span.End()

	return d.wasteRepo.AddWasteToUser(ctxTrace, userID, waste)
}
