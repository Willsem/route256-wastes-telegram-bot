package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

//go:generate mockery --name=telegramClient --dir . --output ./mocks --exported
type telegramClient interface {
	SendMessage(userID int64, text string) error
	SendMessageWithoutRemovingKeyboard(userID int64, text string) error
	SendKeyboard(userID int64, text string, rows [][]string) error
	GetUpdatesChan() chan *models.Message
}

//go:generate mockery --name=wasteRepository --dir . --output ./mocks --exported
type wasteRepository interface {
	GetReportLastWeek(userID int64) ([]models.CategoryReport, error)
	GetReportLastMonth(userID int64) ([]models.CategoryReport, error)
	GetReportLastYear(userID int64) ([]models.CategoryReport, error)

	AddWasteToUser(userID int64, waste *models.Waste) error
}

//go:generate mockery --name=exchangeRepository --dir . --output ./mocks --exported
type exchangesRepository interface {
	GetDefaultCurrency() string
	GetUsedCurrencies() []string
	GetExchange(currency string) (float64, error)
	GetDesignation(currency string) (string, error)
}

type messageContext int

const (
	noContext messageContext = iota
	addWaste
	changeCurrency
)

type reportPeriod string

const (
	week  reportPeriod = "week"
	month reportPeriod = "month"
	year  reportPeriod = "year"
)

const (
	messageWasteNotFound            = "Траты за указанный период не найдены"
	messageInternalError            = "Внутренняя ошибка"
	messageIncorrectFormat          = "Неправильный формат"
	messageSuccessfulAddWaste       = "Трата успешно добавлена"
	messageIncorrectContext         = "Неизвестное состояние пользователя, состояние сброшено до стандартного"
	messageChooseCurrency           = "Выберите валюту из предложенных на клавиатуре"
	messageSuccessfulChangeCurrency = "Валюта успешно изменена на "
	messageDefaultCurrency          = "Выбрана валюта по умолчанию "

	messageAddWaste = `Для добавления траты введите сообщение в формате:

<Название категории>
<Сумма траты>
<Дата траты в формате DD.MM.YYYY> (необязательно)`

	messageHelp = `Данный бот предназначен для ведения трат по категориям

/add - для добавления новой траты
/week - отчет по тратам за последнюю неделю
/month - отчет по тратам за последний месяц
/year - отчет по тратам за последний год
/currency - сменить валюту`
)

const userDateLayout = "02.01.2006"

const convertToMainCurrency = 100.0

type Bot struct {
	tgClient     telegramClient
	wasteRepo    wasteRepository
	exchangeRepo exchangesRepository
	logger       log.Logger

	userContext  map[int64]messageContext
	userCurrency map[int64]string

	cancel context.CancelFunc
	done   chan struct{}
}

func NewBot(
	tgClient telegramClient, wasteRepo wasteRepository, exchangeRepo exchangesRepository, logger log.Logger,
) *Bot {
	return &Bot{
		tgClient:     tgClient,
		wasteRepo:    wasteRepo,
		exchangeRepo: exchangeRepo,
		logger:       logger.With(log.ComponentKey, "Bot"),

		userContext:  make(map[int64]messageContext, 0),
		userCurrency: make(map[int64]string, 0),
	}
}

func (b *Bot) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	b.cancel = cancel
	b.done = make(chan struct{})

	go b.run(ctx)

	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.cancel()

	select {
	case <-b.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Bot) run(ctx context.Context) {
	for {
		select {
		case message := <-b.tgClient.GetUpdatesChan():
			b.logger.With("message", message).Debug("get the message by bot")

			if _, ok := b.userContext[message.From.ID]; !ok {
				b.userContext[message.From.ID] = noContext
			}
			if _, ok := b.userCurrency[message.From.ID]; !ok {
				defaultCurrency := b.exchangeRepo.GetDefaultCurrency()
				b.userCurrency[message.From.ID] = defaultCurrency
				b.sendMessage(message.From.ID, messageDefaultCurrency+defaultCurrency)
			}

			switch message.Text {
			case "/week":
				b.userContext[message.From.ID] = noContext
				b.generateReportForUser(message, week)
			case "/month":
				b.userContext[message.From.ID] = noContext
				b.generateReportForUser(message, month)
			case "/year":
				b.userContext[message.From.ID] = noContext
				b.generateReportForUser(message, year)
			case "/add":
				b.userContext[message.From.ID] = addWaste
				b.sendMessage(message.From.ID, messageAddWaste)
			case "/currency":
				b.userContext[message.From.ID] = changeCurrency
				b.sendChangingCurrencyKeyboard(message.From.ID)
			default:
				b.workWithMessage(message)
			}

		case <-ctx.Done():
			b.logger.WithError(ctx.Err()).Info("bot has been stopped")
			close(b.done)

			return
		}
	}
}

func (b *Bot) generateReportForUser(message *models.Message, period reportPeriod) {
	var (
		report []models.CategoryReport
		err    error
	)

	switch period {
	case week:
		report, err = b.wasteRepo.GetReportLastWeek(message.From.ID)
	case month:
		report, err = b.wasteRepo.GetReportLastMonth(message.From.ID)
	case year:
		report, err = b.wasteRepo.GetReportLastYear(message.From.ID)
	default:
		b.logger.Errorf("unexpected type of period: %s", string(period))
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	if err != nil {
		b.logger.
			WithError(err).
			Warn("failed to generate report during last week for user %d", message.From.ID)

		b.sendMessage(message.From.ID, messageWasteNotFound)
		return
	}

	exchange, err := b.exchangeRepo.GetExchange(b.userCurrency[message.From.ID])
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to get exchange for the currency ", b.userCurrency[message.From.ID])

		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	designation, err := b.exchangeRepo.GetDesignation(b.userCurrency[message.From.ID])
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to get designation for the currency ", b.userCurrency[message.From.ID])

		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	b.sendMessage(message.From.ID, b.generateStringReport(report, period, exchange, designation))
}

func (b *Bot) generateStringReport(
	report []models.CategoryReport, period reportPeriod, currencyExchange float64, currencyDesignation string,
) string {
	textMessageHeader := "Отчет по тратам за "

	switch period {
	case week:
		textMessageHeader += "последнюю неделю:\n\n```\n"
	case month:
		textMessageHeader += "последний месяц:\n\n```\n"
	case year:
		textMessageHeader += "последний год:\n\n```\n"
	default:
		b.logger.Errorf("unexpected type of period: %s", string(period))
		return messageInternalError
	}

	data := make([][]string, 0)
	for _, category := range report {
		data = append(data, []string{
			category.Category,
			fmt.Sprintf("%.2f %s",
				float64(category.Sum)*currencyExchange/convertToMainCurrency, currencyDesignation),
		})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader([]string{"КАТЕГОРИЯ", "ПОТРАЧЕНО"})
	table.SetBorder(false)
	table.AppendBulk(data)

	table.Render()

	return textMessageHeader + tableString.String() + "```"
}

func (b *Bot) workWithMessage(message *models.Message) {
	if _, ok := b.userContext[message.From.ID]; !ok {
		b.userContext[message.From.ID] = noContext
	}

	switch b.userContext[message.From.ID] {
	case noContext:
		b.sendMessage(message.From.ID, messageHelp)
	case addWaste:
		b.addWaste(message)
	case changeCurrency:
		b.changeCurrency(message)
	default:
		b.userContext[message.From.ID] = noContext
		b.sendMessage(message.From.ID, messageIncorrectContext)
	}
}

func (b *Bot) addWaste(message *models.Message) {
	text := message.Text

	lines := strings.Split(text, "\n")

	if len(lines) < 2 || len(lines) > 3 {
		b.sendMessage(message.From.ID, messageIncorrectFormat)
		return
	}

	cost, err := strconv.ParseFloat(lines[1], 64)
	if err != nil {
		b.sendMessage(message.From.ID, messageIncorrectFormat)
		return
	}

	date := message.Date
	if len(lines) == 3 {
		date, err = time.Parse(userDateLayout, lines[2])
		if err != nil {
			b.sendMessage(message.From.ID, messageIncorrectFormat)
			return
		}
	}

	exchange, err := b.exchangeRepo.GetExchange(b.userCurrency[message.From.ID])
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to get exchange", b.userCurrency[message.From.ID])
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	waste := models.NewWaste(lines[0], int64(cost/exchange*convertToMainCurrency), date)
	err = b.wasteRepo.AddWasteToUser(message.From.ID, waste)
	if err != nil {
		b.sendMessage(message.From.ID, messageInternalError)
		b.logger.WithError(err).Error("failed to add waste")
		return
	}

	b.userContext[message.From.ID] = noContext
	b.sendMessage(message.From.ID, messageSuccessfulAddWaste)
}

func (b *Bot) changeCurrency(message *models.Message) {
	_, err := b.exchangeRepo.GetExchange(message.Text)
	if err != nil {
		b.sendMessageWithoutRemovingKeyboard(message.From.ID, messageChooseCurrency)
		return
	}

	b.userContext[message.From.ID] = noContext
	b.userCurrency[message.From.ID] = message.Text
	b.sendMessage(message.From.ID, messageSuccessfulChangeCurrency+message.Text)
}

func (b *Bot) sendChangingCurrencyKeyboard(userID int64) {
	defaultCurrency := b.exchangeRepo.GetDefaultCurrency()
	usedCurrencies := b.exchangeRepo.GetUsedCurrencies()

	usedCurrenciesKeyboardButtons := make([][]string, len(usedCurrencies)+1)
	usedCurrenciesKeyboardButtons[0] = []string{defaultCurrency}
	for i := range usedCurrencies {
		usedCurrenciesKeyboardButtons[i+1] = []string{usedCurrencies[i]}
	}
	b.sendKeyboard(userID, messageChooseCurrency, usedCurrenciesKeyboardButtons)
}

func (b *Bot) sendKeyboard(userID int64, text string, buttons [][]string) {
	err := b.tgClient.SendKeyboard(userID, text, buttons)
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to send keyboard")
	}
}

func (b *Bot) sendMessage(userID int64, text string) {
	err := b.tgClient.SendMessage(userID, text)
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to send message")
	}
}

func (b *Bot) sendMessageWithoutRemovingKeyboard(userID int64, text string) {
	err := b.tgClient.SendMessageWithoutRemovingKeyboard(userID, text)
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to send message")
	}
}
