package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

const warningLimitCoeff = 0.9

const userDateLayout = "02.01.2006"

const (
	messageAddResponse = `Для добавления траты введите сообщение в формате:

<Название категории>
<Сумма траты>
<Дата траты в формате DD.MM.YYYY> (необязательно)`

	messageSuccessfulAddWaste = "Трата успешно добавлена"
	messageWarningLimit       = "До превышения лимита за текущий месяц осталось:"
	messageLimitExceeded      = "Лимит на текущий месяц превышен на"
)

func (h *MessageHandlers) addHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	err := h.userContextService.SetContext(ctx, message.From.ID, enums.AddWaste)
	if err != nil {
		return nil, fmt.Errorf("failed to set user context: %w", err)
	}

	return &bot.MessageResponse{
		Message: messageAddResponse,
	}, nil
}

func (h *MessageHandlers) addWaste(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	text := message.Text
	lines := strings.Split(text, "\n")

	if len(lines) < 2 || len(lines) > 3 {
		return &bot.MessageResponse{
			Message: messageIncorrectFormat,
		}, nil
	}

	cost, err := strconv.ParseFloat(lines[1], 64)
	if err != nil {
		return &bot.MessageResponse{
			Message: messageIncorrectFormat,
		}, nil
	}

	date := message.Date
	if len(lines) == 3 {
		date, err = time.Parse(userDateLayout, lines[2])
		if err != nil {
			return &bot.MessageResponse{
				Message: messageIncorrectFormat,
			}, nil
		}
	}

	exchange, designation, err := h.getExchangeOfUser(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchage and designation for user: %w", err)
	}

	waste := models.NewWaste(lines[0], int64(cost/exchange*convertToMainCurrency), date)
	_, err = h.wasteRepo.AddWasteToUser(ctx, message.From.ID, waste)
	if err != nil {
		return nil, fmt.Errorf("failed to add waste: %w", err)
	}

	err = h.userContextService.SetContext(ctx, message.From.ID, enums.NoContext)
	if err != nil {
		return nil, fmt.Errorf("failed to set context for user: %w", err)
	}

	msg := messageSuccessfulAddWaste + "\n"

	sum, err := h.wasteRepo.SumOfWastesAfterDate(ctx, message.From.ID, getFirstDayOfMonth())
	if err != nil {
		return nil, fmt.Errorf("failed to get sum of wastes: %w", err)
	}

	limit, err := h.userRepo.GetWasteLimit(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get limit of wastes: %w", err)
	}

	if limit != nil {
		fsum := float64(sum)
		flimit := float64(*limit)

		if fsum > flimit {
			diff := h.convertFromDefaultCurrency(uint64(sum)-(*limit), exchange)
			msg += fmt.Sprintf("%s %.2f %s", messageLimitExceeded, diff, designation)
		} else if fsum > warningLimitCoeff*flimit {
			diff := h.convertFromDefaultCurrency((*limit)-uint64(sum), exchange)
			msg += fmt.Sprintf("%s %.2f %s", messageWarningLimit, diff, designation)
		}
	}

	return &bot.MessageResponse{
		Message: msg,
	}, nil
}

func getFirstDayOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}
