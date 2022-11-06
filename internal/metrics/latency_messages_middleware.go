package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

func LatencyMetricMiddleware(commands []string) bot.MessageMiddleware {
	latencyMessages := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "telegram_messages_response_latency",
			Help:    "Duration of message response",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0},
		}, []string{"type"},
	)

	middleware := func(next bot.MessageHandler) bot.MessageHandler {
		return func(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
			startTime := time.Now()
			response, err := next(ctx, message)
			duration := time.Since(startTime)

			latencyMessages.WithLabelValues(messageType(message.Text, commands)).Observe(duration.Seconds())

			return response, err
		}
	}

	return middleware
}
