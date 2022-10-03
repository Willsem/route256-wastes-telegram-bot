package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

//go:generate mockery --name=telegramClient --dir .  --output ./mocks --exported
type telegramClient interface {
	SendMessage(userID int64, text string) error
	GetUpdatesChan() chan *models.Message
}

//go:generate mockery --name=wasteRepository --dir .  --output ./mocks --exported
type wasteRepository interface {
	GetReportLastWeek(userID int64) ([]models.CategoryReport, error)
	GetReportLastMonth(userID int64) ([]models.CategoryReport, error)
	GetReportLastYear(userID int64) ([]models.CategoryReport, error)

	AddWasteToUser(userID int64, waste *models.Waste) error
}

type messageContext int

const (
	noContext messageContext = iota
	addWaste
)

type reportPeriod string

const (
	week  reportPeriod = "week"
	month reportPeriod = "month"
	year  reportPeriod = "year"
)

const (
	messageWasteNotFound       = "Траты за указанный период не найдены"
	messageInternalError       = "Внутренняя ошибка"
	messageIncorrectFormat     = "Неправильный формат"
	messageSuccessfullAddWaste = "Трата успешно добавлена"

	messageAddWaste = `Для добавления траты введите сообщение в формате:

<Название категории>
<Сумма траты>
<Дата траты в формате DD.MM.YYYY> (необязательно)`

	messageHelp = `Данный бот предназначен для ведения трат по категориям

/add - для добавления новой траты
/week - отчет по тратам за последнюю неделю
/month - отчет по тратам за последний месяц
/year - отчет по тратам за последний год`
)

const userDateLayout = "02.01.2006"

type Bot struct {
	tgClient  telegramClient
	wasteRepo wasteRepository
	logger    log.Logger

	userContext map[int64]messageContext
}

func NewBot(tgClient telegramClient, wasteRepo wasteRepository, logger log.Logger) *Bot {
	return &Bot{
		tgClient:    tgClient,
		wasteRepo:   wasteRepo,
		logger:      logger,
		userContext: make(map[int64]messageContext, 0),
	}
}

func (b *Bot) Run() {
	for message := range b.tgClient.GetUpdatesChan() {
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
		default:
			b.workWithMessage(message)
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
			Error("failed to generate report during last week for user %d", message.From.ID)

		b.sendMessage(message.From.ID, messageWasteNotFound)
		return
	}

	b.sendMessage(message.From.ID, b.generateStringReport(report, period))
}

func (b *Bot) generateStringReport(report []models.CategoryReport, period reportPeriod) string {
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
			strconv.Itoa(category.Sum),
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
	}
}

func (b *Bot) addWaste(message *models.Message) {
	text := message.Text

	lines := strings.Split(text, "\n")

	if len(lines) < 2 || len(lines) > 3 {
		b.sendMessage(message.From.ID, messageIncorrectFormat)
		return
	}

	cost, err := strconv.Atoi(lines[1])
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

	waste := models.NewWaste(lines[0], cost, date)
	err = b.wasteRepo.AddWasteToUser(message.From.ID, waste)
	if err != nil {
		b.sendMessage(message.From.ID, messageInternalError)
		b.logger.WithError(err).Error("failed to add waste")
		return
	}

	b.sendMessage(message.From.ID, messageSuccessfullAddWaste)
}

func (b *Bot) sendMessage(userID int64, text string) {
	err := b.tgClient.SendMessage(userID, text)
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to send message")
	}
}
