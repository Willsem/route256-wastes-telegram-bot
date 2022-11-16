package handlers

import (
	"context"
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

//go:generate mockery --name=userRepository --dir . --output ./mocks --exported
type userRepository interface {
	UserExists(ctx context.Context, id int64) (bool, error)
	AddUser(ctx context.Context, user *models.User) (*models.User, error)

	SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error)
	GetWasteLimit(ctx context.Context, id int64) (*uint64, error)
}

//go:generate mockery --name=wasteRepository --dir . --output ./mocks --exported
type wasteRepository interface {
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

//go:generate mockery --name=userContextService --dir . --output ./mocks --exported
type userContextService interface {
	SetContext(ctx context.Context, userID int64, context enums.UserContext) error
	GetContext(ctx context.Context, userID int64) (enums.UserContext, error)
	SetCurrency(ctx context.Context, userID int64, currency string) error
	GetCurrency(ctx context.Context, userID int64) (string, error)
}

//go:generate mockery --name=kafkaProducer --dir . --output ./mocks --exported
type kafkaProducer interface {
	SendMessage(ctx context.Context, key []byte, value []byte) error
}

type MessageHandlers struct {
	userRepo           userRepository
	wasteRepo          wasteRepository
	exchangeService    exchangeService
	userContextService userContextService
	kafkaProducer      kafkaProducer
}

func NewMessageHandlers(
	userRepo userRepository,
	wasteRepo wasteRepository,
	exchangeService exchangeService,
	userContextService userContextService,
	kafkaProducer kafkaProducer,
) *MessageHandlers {
	return &MessageHandlers{
		userRepo:           userRepo,
		wasteRepo:          wasteRepo,
		exchangeService:    exchangeService,
		userContextService: userContextService,
		kafkaProducer:      kafkaProducer,
	}
}

func (h *MessageHandlers) GetHandlers() map[string]bot.MessageHandler {
	return map[string]bot.MessageHandler{
		"/add":      h.addHandler,
		"/setLimit": h.setLimitHandler,
		"/getLimit": h.getLimitHandler,
		"/week":     h.weekHandler,
		"/month":    h.monthHandler,
		"/year":     h.yearHandler,
		"/currency": h.currencyHandler,
		"default":   h.defaultHandler,
	}
}
