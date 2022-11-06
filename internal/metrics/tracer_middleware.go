package metrics

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TracingMiddleware(tracerProvider *tracesdk.TracerProvider, commands []string) bot.MessageMiddleware {
	tracer := tracerProvider.Tracer("message-middleware")

	middleware := func(next bot.MessageHandler) bot.MessageHandler {
		return func(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
			ctxTrace, span := tracer.Start(ctx, messageType(message.Text, commands))
			defer span.End()

			return next(ctxTrace, message)
		}
	}

	return middleware
}
