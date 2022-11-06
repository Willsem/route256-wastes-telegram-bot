package handlers

import (
	"context"
	"fmt"
)

const (
	messageIncorrectFormat = "Неправильный формат"

	convertToMainCurrency = 100.0
)

func (h *MessageHandlers) convertFromDefaultCurrency(money uint64, exchange float64) float64 {
	return float64(money) * exchange / convertToMainCurrency
}

func (h *MessageHandlers) convertToDefaultCurrency(money float64, exchange float64) uint64 {
	return uint64(float64(money) / exchange * convertToMainCurrency)
}

func (h *MessageHandlers) getExchangeOfUser(ctx context.Context, userID int64) (float64, string, error) {
	currency, err := h.userContextService.GetCurrency(ctx, userID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get user currency: %w", err)
	}

	exchange, err := h.exchangeService.GetExchange(currency)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get exchange of user: %w", err)
	}

	designation, err := h.exchangeService.GetDesignation(currency)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get designation os user: %w", err)
	}

	return exchange, designation, nil
}
