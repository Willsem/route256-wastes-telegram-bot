package metrics

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type TelegramClientTracerDecorator struct {
	tgClient telegramClient
	tracer   trace.Tracer
}

func NewTelegramClientTracerDecorator(tgClient telegramClient, tracerProvider *tracesdk.TracerProvider) *TelegramClientTracerDecorator {
	return &TelegramClientTracerDecorator{
		tgClient: tgClient,
		tracer:   tracerProvider.Tracer("telegram-bot-client"),
	}
}

func (d *TelegramClientTracerDecorator) SendMessage(ctx context.Context, userID int64, text string) error {
	ctxTrace, span := d.tracer.Start(ctx, "SendMessage")
	defer span.End()

	return d.tgClient.SendMessage(ctxTrace, userID, text)
}

func (d *TelegramClientTracerDecorator) SendMessageWithoutRemovingKeyboard(ctx context.Context, userID int64, text string) error {
	ctxTrace, span := d.tracer.Start(ctx, "SendMessageWithoutRemovingKeyboard")
	defer span.End()

	return d.tgClient.SendMessageWithoutRemovingKeyboard(ctxTrace, userID, text)
}

func (d *TelegramClientTracerDecorator) SendKeyboard(ctx context.Context, userID int64, text string, rows [][]string) error {
	ctxTrace, span := d.tracer.Start(ctx, "SendKeyboard")
	defer span.End()

	return d.tgClient.SendKeyboard(ctxTrace, userID, text, rows)
}

func (d *TelegramClientTracerDecorator) GetUpdatesChan() <-chan *models.Message {
	return d.tgClient.GetUpdatesChan()
}
