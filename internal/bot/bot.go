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

//go:generate mockery --name=userRepository --dir . --output ./mocks --exported
type userRepository interface {
	UserExists(ctx context.Context, id int64) (bool, error)
	AddUser(ctx context.Context, user *models.User) (*models.User, error)

	SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error)
	GetWasteLimit(ctx context.Context, id int64) (*uint64, error)
}

//go:generate mockery --name=wasteRepository --dir . --output ./mocks --exported
type wasteRepository interface {
	GetReportLastWeek(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	GetReportLastMonth(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	GetReportLastYear(ctx context.Context, userID int64) ([]*models.CategoryReport, error)
	SumOfWastesAfterDate(ctx context.Context, userID int64, date time.Time) (int64, error)

	AddWasteToUser(ctx context.Context, userID int64, waste *models.Waste) (*models.Waste, error)
}

//go:generate mockery --name=exchangeService --dir . --output ./mocks --exported
type exchangeService interface {
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
	setLimit
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
	messageWarningLimit             = "До превышения лимита за текущий месяц осталось:"
	messageLimitExceeded            = "Лимит на текущий месяц превышен на"
	messageSetLimit                 = "Введите желаемый лимит в виде вещественного числа в текущей валюте"
	messageGetLimit                 = "Текущий лимит на месяц:"
	messageNullLimit                = "Лимит на месяц не установлен"
	messageSuccessfulSetLimit       = "Лимит трат за месяц успешно установлен"

	messageAddWaste = `Для добавления траты введите сообщение в формате:

<Название категории>
<Сумма траты>
<Дата траты в формате DD.MM.YYYY> (необязательно)`

	messageHelp = `Данный бот предназначен для ведения трат по категориям

/add - для добавления новой траты
/setLimit - установить лимит на месяц
/getLimit - узнать текущий лимит на месяц
/week - отчет по тратам за последнюю неделю
/month - отчет по тратам за последний месяц
/year - отчет по тратам за последний год
/currency - сменить валюту`
)

const userDateLayout = "02.01.2006"

const (
	convertToMainCurrency = 100.0
	warningLimitCoeff     = 0.9
)

type Bot struct {
	tgClient telegramClient

	wasteRepo wasteRepository
	userRepo  userRepository

	exchangeService exchangeService
	logger          log.Logger

	userContext  map[int64]messageContext
	userCurrency map[int64]string

	cancel context.CancelFunc
	done   chan struct{}
}

func NewBot(
	tgClient telegramClient,
	userRepo userRepository, wasteRepo wasteRepository,
	exchangeService exchangeService, logger log.Logger,
) *Bot {
	return &Bot{
		tgClient: tgClient,

		wasteRepo: wasteRepo,
		userRepo:  userRepo,

		exchangeService: exchangeService,
		logger:          logger.With(log.ComponentKey, "Bot"),

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
				defaultCurrency := b.exchangeService.GetDefaultCurrency()
				b.userCurrency[message.From.ID] = defaultCurrency
				b.sendMessage(message.From.ID, messageDefaultCurrency+defaultCurrency)
			}

			userExists, err := b.userRepo.UserExists(ctx, message.From.ID)
			if err != nil {
				b.logger.WithError(err).
					Error("failed to check that user exists")
				b.sendMessage(message.From.ID, messageInternalError)
				break
			}

			if !userExists {
				user := message.From
				_, err = b.userRepo.AddUser(ctx, models.NewUser(
					user.ID, user.FirstName, user.LastName, user.UserName))
				if err != nil {
					b.logger.WithError(err).
						Error("failed to create new user")
					b.sendMessage(message.From.ID, messageInternalError)
					break
				}
			}

			switch message.Text {
			case "/week":
				b.userContext[message.From.ID] = noContext
				b.generateReportForUser(ctx, message, week)
			case "/month":
				b.userContext[message.From.ID] = noContext
				b.generateReportForUser(ctx, message, month)
			case "/year":
				b.userContext[message.From.ID] = noContext
				b.generateReportForUser(ctx, message, year)
			case "/add":
				b.userContext[message.From.ID] = addWaste
				b.sendMessage(message.From.ID, messageAddWaste)
			case "/currency":
				b.userContext[message.From.ID] = changeCurrency
				b.sendChangingCurrencyKeyboard(message.From.ID)
			case "/setLimit":
				b.userContext[message.From.ID] = setLimit
				b.sendMessage(message.From.ID, messageSetLimit)
			case "/getLimit":
				b.userContext[message.From.ID] = noContext
				b.getLimit(ctx, message)
			default:
				b.workWithMessage(ctx, message)
			}

		case <-ctx.Done():
			b.logger.WithError(ctx.Err()).Info("bot has been stopped")
			close(b.done)

			return
		}
	}
}

func (b *Bot) generateReportForUser(ctx context.Context, message *models.Message, period reportPeriod) {
	var (
		report []*models.CategoryReport
		err    error
	)

	switch period {
	case week:
		report, err = b.wasteRepo.GetReportLastWeek(ctx, message.From.ID)
	case month:
		report, err = b.wasteRepo.GetReportLastMonth(ctx, message.From.ID)
	case year:
		report, err = b.wasteRepo.GetReportLastYear(ctx, message.From.ID)
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

	if len(report) == 0 {
		b.sendMessage(message.From.ID, messageWasteNotFound)
		return
	}

	exchange, designation, err := b.getExchangeOfUser(message.From.ID)
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to get exchange and designation for the currency ", b.userCurrency[message.From.ID])

		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	b.sendMessage(message.From.ID, b.generateStringReport(report, period, exchange, designation))
}

func (b *Bot) generateStringReport(
	report []*models.CategoryReport, period reportPeriod, currencyExchange float64, currencyDesignation string,
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

	return textMessageHeader + tableString.String() + "```"
}

func (b *Bot) workWithMessage(ctx context.Context, message *models.Message) {
	if _, ok := b.userContext[message.From.ID]; !ok {
		b.userContext[message.From.ID] = noContext
	}

	switch b.userContext[message.From.ID] {
	case noContext:
		b.sendMessage(message.From.ID, messageHelp)
	case addWaste:
		b.addWaste(ctx, message)
	case changeCurrency:
		b.changeCurrency(message)
	case setLimit:
		b.setLimit(ctx, message)
	default:
		b.userContext[message.From.ID] = noContext
		b.sendMessage(message.From.ID, messageIncorrectContext)
	}
}

func (b *Bot) addWaste(ctx context.Context, message *models.Message) {
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

	exchange, err := b.exchangeService.GetExchange(b.userCurrency[message.From.ID])
	if err != nil {
		b.logger.
			WithError(err).
			Error("failed to get exchange", b.userCurrency[message.From.ID])
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	waste := models.NewWaste(lines[0], int64(cost/exchange*convertToMainCurrency), date)
	_, err = b.wasteRepo.AddWasteToUser(ctx, message.From.ID, waste)
	if err != nil {
		b.sendMessage(message.From.ID, messageInternalError)
		b.logger.WithError(err).Error("failed to add waste")
		return
	}

	b.userContext[message.From.ID] = noContext
	b.sendMessage(message.From.ID, messageSuccessfulAddWaste)

	sum, err := b.wasteRepo.SumOfWastesAfterDate(ctx, message.From.ID, getFirstDayOfMonth())
	if err != nil {
		b.logger.WithError(err).Error("failed to calculate sum of wastes")
		return
	}

	limit, err := b.userRepo.GetWasteLimit(ctx, message.From.ID)
	if err != nil {
		b.logger.WithError(err).Error("failed to get the wasting limit")
		return
	}

	if limit != nil {
		msg := ""

		fsum := float64(sum)
		flimit := float64(*limit)

		exchange, designation, err := b.getExchangeOfUser(message.From.ID)
		if err != nil {
			b.logger.WithError(err).Errorf("failed to get exchange %s for the user",
				b.userCurrency[message.From.ID])
			return
		}

		if fsum > flimit {
			diff := (fsum - flimit) * exchange / convertToMainCurrency
			msg = fmt.Sprintf("%s %.2f %s", messageLimitExceeded, diff, designation)
		} else if fsum > warningLimitCoeff*flimit {
			diff := (flimit - fsum) * exchange / convertToMainCurrency
			msg = fmt.Sprintf("%s %.2f %s", messageWarningLimit, diff, designation)
		}

		if msg != "" {
			b.sendMessage(message.From.ID, msg)
		}
	}
}

func (b *Bot) changeCurrency(message *models.Message) {
	_, err := b.exchangeService.GetExchange(message.Text)
	if err != nil {
		b.sendMessageWithoutRemovingKeyboard(message.From.ID, messageChooseCurrency)
		return
	}

	b.userContext[message.From.ID] = noContext
	b.userCurrency[message.From.ID] = message.Text
	b.sendMessage(message.From.ID, messageSuccessfulChangeCurrency+message.Text)
}

func (b *Bot) setLimit(ctx context.Context, message *models.Message) {
	limit, err := strconv.ParseFloat(message.Text, 64)
	if err != nil {
		b.sendMessage(message.From.ID, messageIncorrectFormat)
		return
	}

	exchange, _, err := b.getExchangeOfUser(message.From.ID)
	if err != nil {
		b.logger.WithError(err).Error("failed to get exchange of user")
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	_, err = b.userRepo.SetWasteLimit(ctx, message.From.ID,
		uint64(limit*convertToMainCurrency/exchange))
	if err != nil {
		b.logger.WithError(err).Error("failed to set waste limit of user")
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	b.sendMessage(message.From.ID, messageSuccessfulSetLimit)
}

func (b *Bot) getLimit(ctx context.Context, message *models.Message) {
	limit, err := b.userRepo.GetWasteLimit(ctx, message.From.ID)
	if err != nil {
		b.logger.WithError(err).Error("failed to get limit for user")
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	if limit == nil {
		b.sendMessage(message.From.ID, messageNullLimit)
		return
	}

	exchange, designation, err := b.getExchangeOfUser(message.From.ID)
	if err != nil {
		b.logger.WithError(err).Error("failed to get exchange for user")
		b.sendMessage(message.From.ID, messageInternalError)
		return
	}

	b.sendMessage(message.From.ID, fmt.Sprintf("%s %.2f %s",
		messageGetLimit, float64(*limit)*exchange/convertToMainCurrency, designation))
}

func (b *Bot) sendChangingCurrencyKeyboard(userID int64) {
	defaultCurrency := b.exchangeService.GetDefaultCurrency()
	usedCurrencies := b.exchangeService.GetUsedCurrencies()

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

func (b *Bot) getExchangeOfUser(userID int64) (float64, string, error) {
	exchange, err := b.exchangeService.GetExchange(b.userCurrency[userID])
	if err != nil {
		return 0, "", err
	}

	designation, err := b.exchangeService.GetDesignation(b.userCurrency[userID])
	if err != nil {
		return 0, "", err
	}

	return exchange, designation, nil
}

func getFirstDayOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}
