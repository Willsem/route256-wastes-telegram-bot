package handlers

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

const (
	messageSetLimitResponse   = "Введите желаемый лимит в виде положительного вещественного числа в текущей валюте"
	messageSuccessfulSetLimit = "Лимит трат за месяц успешно установлен"
)

func (h *MessageHandlers) setLimitHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	err := h.userContextService.SetContext(ctx, message.From.ID, enums.SetLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to set context for user: %w", err)
	}

	return &bot.MessageResponse{
		Message: messageSetLimitResponse,
	}, nil
}

func (h *MessageHandlers) setLimit(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	limit, err := strconv.ParseFloat(message.Text, 64)
	if err != nil {
		return &bot.MessageResponse{
			Message: messageIncorrectFormat,
		}, nil
	}

	if limit < 0 {
		return &bot.MessageResponse{
			Message: messageIncorrectFormat,
		}, nil
	}

	exchange, _, err := h.getExchangeOfUser(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange and designation for user: %w", err)
	}

	_, err = h.userRepo.SetWasteLimit(ctx, message.From.ID, h.convertToDefaultCurrency(limit, exchange))
	if err != nil {
		return nil, fmt.Errorf("failed to set waste for user: %w", err)
	}

	err = h.userContextService.SetContext(ctx, message.From.ID, enums.NoContext)
	if err != nil {
		return nil, fmt.Errorf("failed to set context for user: %w", err)
	}

	return &bot.MessageResponse{
		Message: messageSuccessfulSetLimit,
	}, nil
}
