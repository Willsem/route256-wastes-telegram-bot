package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

func AmountMetricMiddleware(commands []string) bot.MessageMiddleware {
	countMessages := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "telegram_messages_count",
			Help: "Count of messages",
		}, []string{"type"},
	)

	middleware := func(next bot.MessageHandler) bot.MessageHandler {
		return func(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
			countMessages.WithLabelValues(messageType(message.Text, commands)).Inc()
			return next(ctx, message)
		}
	}

	return middleware
}
