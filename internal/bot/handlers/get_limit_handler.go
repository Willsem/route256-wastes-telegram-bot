package handlers

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

const (
	messageGetLimit  = "Текущий лимит на месяц:"
	messageNullLimit = "Лимит на месяц не установлен"
)

func (h *MessageHandlers) getLimitHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	err := h.userContextService.SetContext(ctx, message.From.ID, enums.AddWaste)
	if err != nil {
		return nil, fmt.Errorf("failed to set user context: %w", err)
	}

	limit, err := h.userRepo.GetWasteLimit(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the limit: %w", err)
	}

	if limit == nil {
		return &bot.MessageResponse{
			Message: messageNullLimit,
		}, nil
	}

	exchange, designation, err := h.getExchangeOfUser(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange and designation of user: %w", err)
	}

	return &bot.MessageResponse{
		Message: fmt.Sprintf("%s %.2f %s",
			messageGetLimit, h.convertFromDefaultCurrency(*limit, exchange), designation),
	}, nil
}
