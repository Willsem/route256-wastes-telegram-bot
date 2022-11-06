package handlers

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

const (
	messageHelp = `**Данный бот предназначен для ведения трат по категориям**

/add - для добавления новой траты
/setLimit - установить лимит на месяц
/getLimit - узнать текущий лимит на месяц
/week - отчет по тратам за последнюю неделю
/month - отчет по тратам за последний месяц
/year - отчет по тратам за последний год
/currency - сменить валюту`

	messageIncorrectContext = "Неизвестное состояние пользователя, состояние сброшено до стандартного"
)

func (h *MessageHandlers) defaultHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	userContext, err := h.userContextService.GetContext(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	switch userContext {
	case enums.NoContext:
		return &bot.MessageResponse{
			Message: messageHelp,
		}, nil

	case enums.AddWaste:
		return h.addWaste(ctx, message)

	case enums.ChangeCurrency:
		return h.changeCurrency(ctx, message)

	case enums.SetLimit:
		return h.setLimit(ctx, message)

	default:
		err := h.userContextService.SetContext(ctx, message.From.ID, enums.NoContext)
		if err != nil {
			return nil, fmt.Errorf("failed to set user context: %w", err)
		}
		return &bot.MessageResponse{
			Message: messageIncorrectContext,
		}, nil
	}
}
