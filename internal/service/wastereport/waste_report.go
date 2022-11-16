package wastereport

import (
	"context"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/requests"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

const convertToMainCurrency = 100.0
const messageWasteNotFound = "Траты за указанный период не найдены"

//go:generate mockery --name=wasteRepository --dir . --output ./mocks --exported
type wasteRepository interface {
	GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
}

//go:generate mockery --name=consumerMessages --dir . --output ./mocks --exported
type consumerMessages interface {
	GetMessageChan() <-chan *models.KafkaMessage
}

//go:generate mockery --name=telegramClient --dir . --output ./mocks --exported
type telegramClient interface {
	SendMessage(ctx context.Context, userID int64, text string, command enums.CommandType) error
}

type Service struct {
	consumer  consumerMessages
	wasteRepo wasteRepository
	tgClient  telegramClient

	logger log.Logger

	cancel context.CancelFunc
	done   chan struct{}
}

func NewService(consumer consumerMessages, wasteRepo wasteRepository, tgClient telegramClient, logger log.Logger) *Service {
	return &Service{
		consumer:  consumer,
		wasteRepo: wasteRepo,
		tgClient:  tgClient,

		logger: logger.With(log.ComponentKey, "Waste report"),
	}
}

func (s *Service) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	s.cancel = cancel
	s.done = make(chan struct{})

	go s.run(ctx)

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	s.cancel()

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Service) run(ctx context.Context) {
	for {
		select {
		case msg := <-s.consumer.GetMessageChan():
			var req requests.GetReport
			err := req.UnmarshalJSON(msg.Message)
			if err != nil {
				s.logger.
					WithError(err).
					With("recieved message", msg).
					Warn("failed to unmarshall message")
			}
			s.sendReport(ctx, req)

		case <-ctx.Done():
			s.logger.WithError(ctx.Err()).Info("waste report service has been closed")
			close(s.done)

			return

		//nolint:staticcheck // does not recieve the messages without default branch
		default:
		}
	}
}

func (s *Service) sendReport(ctx context.Context, req requests.GetReport) {
	var report []*models.CategoryReport
	var err error
	var command enums.CommandType

	switch req.Period {
	case requests.PeriodWeek:
		report, err = s.wasteRepo.GetReportLastWeek(ctx, req.UserID)
		command = enums.CommandTypeWeekReport
	case requests.PeriodMonth:
		report, err = s.wasteRepo.GetReportLastMonth(ctx, req.UserID)
		command = enums.CommandTypeMonthReport
	case requests.PeriodYear:
		report, err = s.wasteRepo.GetReportLastYear(ctx, req.UserID)
		command = enums.CommandTypeYearReport
	default:
		s.logger.With("report request", req).Warn("unexpected type of period")
		return
	}

	if err != nil {
		s.logger.WithError(err).Error("failed to get the report from repository")
	}

	msg := ""
	if len(report) == 0 {
		msg = messageWasteNotFound
	} else {
		stringReport, err := s.generateStringReport(report, req.Period, req.CurrencyExchange, req.CurrencyDesignation)
		if err != nil {
			s.logger.WithError(err).Error("failed to generate string report")
		}

		msg = stringReport
	}

	err = s.tgClient.SendMessage(ctx, req.UserID, msg, command)
	if err != nil {
		s.logger.WithError(err).Error("failed to send the message")
	}
}

func (s *Service) generateStringReport(
	report []*models.CategoryReport, period requests.Period, currencyExchange float64, currencyDesignation string,
) (string, error) {
	textMessageHeader := "Отчет по тратам за "

	switch period {
	case requests.PeriodWeek:
		textMessageHeader += "последнюю неделю:\n\n```\n"
	case requests.PeriodMonth:
		textMessageHeader += "последний месяц:\n\n```\n"
	case requests.PeriodYear:
		textMessageHeader += "последний год:\n\n```\n"
	default:
		return "", fmt.Errorf("unexpected type of period: %d", period)
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
