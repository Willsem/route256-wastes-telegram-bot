package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

//go:generate mockery --name=telegramClient --dir . --output ./mocks --exported
type telegramClient interface {
	SendMessage(ctx context.Context, userID int64, text string) error
	SendMessageWithoutRemovingKeyboard(ctx context.Context, userID int64, text string) error
	SendKeyboard(ctx context.Context, userID int64, text string, rows [][]string) error
	GetUpdatesChan() chan *models.Message
}

type TelegramClientLatencyDecorator struct {
	tgClient telegramClient
	latency  *prometheus.HistogramVec
}

func NewTelegramClientLatencyDecorator(tgClient telegramClient) *TelegramClientLatencyDecorator {
	return &TelegramClientLatencyDecorator{
		tgClient: tgClient,
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "latency_telegram_client",
			Help:    "Duration of TelegramCLient methods",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0},
		}, []string{"method"}),
	}
}

func (d *TelegramClientLatencyDecorator) SendMessage(ctx context.Context, userID int64, text string) error {
	startTime := time.Now()
	err := d.tgClient.SendMessage(ctx, userID, text)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SendMessage").Observe(duration.Seconds())

	return err
}

func (d *TelegramClientLatencyDecorator) SendMessageWithoutRemovingKeyboard(ctx context.Context, userID int64, text string) error {
	startTime := time.Now()
	err := d.tgClient.SendMessageWithoutRemovingKeyboard(ctx, userID, text)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SendMessageWithoutRemovingKeyboard").Observe(duration.Seconds())

	return err
}

func (d *TelegramClientLatencyDecorator) SendKeyboard(ctx context.Context, userID int64, text string, rows [][]string) error {
	startTime := time.Now()
	err := d.tgClient.SendKeyboard(ctx, userID, text, rows)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SendKeyboard").Observe(duration.Seconds())

	return err
}

func (d *TelegramClientLatencyDecorator) GetUpdatesChan() chan *models.Message {
	return d.tgClient.GetUpdatesChan()
}
