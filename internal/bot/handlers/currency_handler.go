package handlers

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

const (
	messageChooseCurrency           = "Выберите валюту из предложенных на клавиатуре"
	messageSuccessfulChangeCurrency = "Валюта успешно изменена на "
)

func (h *MessageHandlers) currencyHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	err := h.userContextService.SetContext(ctx, message.From.ID, enums.ChangeCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to set context for user: %w", err)
	}

	defaultCurrency := h.exchangeService.GetDefaultCurrency()
	usedCurrencies := h.exchangeService.GetUsedCurrencies()

	usedCurrenciesKeyboardButtons := make([][]string, len(usedCurrencies)+1)
	usedCurrenciesKeyboardButtons[0] = []string{defaultCurrency}
	for i := range usedCurrencies {
		usedCurrenciesKeyboardButtons[i+1] = []string{usedCurrencies[i]}
	}

	return &bot.MessageResponse{
		Message:  messageChooseCurrency,
		Keyboard: usedCurrenciesKeyboardButtons,
	}, nil
}

func (h *MessageHandlers) changeCurrency(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	_, err := h.exchangeService.GetExchange(message.Text)
	if err != nil {
		return &bot.MessageResponse{
			Message:             messageChooseCurrency,
			DoNotRemoveKeyboard: true,
		}, nil
	}

	err = h.userContextService.SetContext(ctx, message.From.ID, enums.NoContext)
	if err != nil {
		return nil, fmt.Errorf("failed to set context for user: %w", err)
	}

	err = h.userContextService.SetCurrency(ctx, message.From.ID, message.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to set currency for user: %w", err)
	}

	return &bot.MessageResponse{
		Message: messageSuccessfulChangeCurrency + message.Text,
	}, nil
}
