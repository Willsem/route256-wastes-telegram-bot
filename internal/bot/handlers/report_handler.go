package handlers

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/requests"
)

const generatingReportMessage = "Отчет генерируется..."

func (h *MessageHandlers) weekHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	return h.generateReportForUser(ctx, message, requests.PeriodWeek)
}

func (h *MessageHandlers) monthHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	return h.generateReportForUser(ctx, message, requests.PeriodMonth)
}

func (h *MessageHandlers) yearHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	return h.generateReportForUser(ctx, message, requests.PeriodYear)
}

func (h *MessageHandlers) generateReportForUser(ctx context.Context, message *models.Message, period requests.Period) (*bot.MessageResponse, error) {
	exchange, designation, err := h.getExchangeOfUser(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchage and designation for the user: %w", err)
	}

	req := requests.GetReport{
		UserID:              message.From.ID,
		Period:              period,
		CurrencyExchange:    exchange,
		CurrencyDesignation: designation,
	}

	key := []byte{}
	value, err := req.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the request: %w", err)
	}

	err = h.kafkaProducer.SendMessage(ctx, key, value)
	if err != nil {
		return nil, fmt.Errorf("failed to send the message to the kafka: %w", err)
	}

	return &bot.MessageResponse{
		Message: generatingReportMessage,
	}, nil
}
