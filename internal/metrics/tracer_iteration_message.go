package metrics

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockery --name=iterationMessage --dir . --output ./mocks --exported
type iterationMessage interface {
	Iterate(ctx context.Context, message *models.Message, handler bot.MessageHandler, logger log.Logger)
}

type IterationMessageTracerDecorator struct {
	iterationMessage iterationMessage
	tracer           trace.Tracer
}

func NewIterationMessageTracerDecorator(iterationMessage iterationMessage, tracerProvider *tracesdk.TracerProvider) *IterationMessageTracerDecorator {
	return &IterationMessageTracerDecorator{
		iterationMessage: iterationMessage,
		tracer:           tracerProvider.Tracer("iteration-message"),
	}
}

func (d *IterationMessageTracerDecorator) Iterate(ctx context.Context, message *models.Message, handler bot.MessageHandler, logger log.Logger) {
	ctxTrace, span := d.tracer.Start(ctx, "IterateMessage")
	defer span.End()

	d.iterationMessage.Iterate(ctxTrace, message, handler, logger)
}
