package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

type reportPeriod string

const (
	week  reportPeriod = "week"
	month reportPeriod = "month"
	year  reportPeriod = "year"
)

const messageWasteNotFound = "Траты за указанный период не найдены"

func (h *MessageHandlers) weekHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	return h.generateReportForUser(ctx, message, week)
}

func (h *MessageHandlers) monthHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	return h.generateReportForUser(ctx, message, month)
}

func (h *MessageHandlers) yearHandler(ctx context.Context, message *models.Message) (*bot.MessageResponse, error) {
	return h.generateReportForUser(ctx, message, year)
}

func (h *MessageHandlers) generateReportForUser(ctx context.Context, message *models.Message, period reportPeriod) (*bot.MessageResponse, error) {
	err := h.userContextService.SetContext(ctx, message.From.ID, enums.AddWaste)
	if err != nil {
		return nil, fmt.Errorf("failed to set user context: %w", err)
	}

	var report []*models.CategoryReport

	switch period {
	case week:
		report, err = h.wasteRepo.GetReportLastWeek(ctx, message.From.ID)
	case month:
		report, err = h.wasteRepo.GetReportLastMonth(ctx, message.From.ID)
	case year:
		report, err = h.wasteRepo.GetReportLastYear(ctx, message.From.ID)
	default:
		return nil, fmt.Errorf("unexpected type of period: %s", period)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get report for the period %s: %w", period, err)
	}

	if len(report) == 0 {
		return &bot.MessageResponse{
			Message: messageWasteNotFound,
		}, nil
	}

	exchange, designation, err := h.getExchangeOfUser(ctx, message.From.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchage and designation for the user: %w", err)
	}

	msg, err := h.generateStringReport(report, period, exchange, designation)
	if err != nil {
		return nil, fmt.Errorf("failed to generate string report: %w", err)
	}

	return &bot.MessageResponse{
		Message: msg,
	}, nil
}

func (h *MessageHandlers) generateStringReport(
	report []*models.CategoryReport, period reportPeriod, currencyExchange float64, currencyDesignation string,
) (string, error) {
	textMessageHeader := "Отчет по тратам за "

	switch period {
	case week:
		textMessageHeader += "последнюю неделю:\n\n```\n"
	case month:
		textMessageHeader += "последний месяц:\n\n```\n"
	case year:
		textMessageHeader += "последний год:\n\n```\n"
	default:
		return "", fmt.Errorf("unexpected type of period: %s", period)
	}

	data := make([][]string, 0)
	sum := 0.0
	for _, category := range report {
		curr := float64(category.Sum) * currencyExchange / convertToMainCurrency
		sum += curr

		data = append(data, []string{
			category.Category,
			fmt.Sprintf("%.2f %s",
				curr, currencyDesignation),
		})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader([]string{"КАТЕГОРИЯ", "ПОТРАЧЕНО"})
	table.SetFooter([]string{"СУММА", fmt.Sprintf("%.2f %s", sum, currencyDesignation)})
	table.AppendBulk(data)

	table.Render()

	return textMessageHeader + tableString.String() + "```", nil
}
