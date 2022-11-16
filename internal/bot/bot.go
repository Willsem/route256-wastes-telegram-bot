package bot

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

const messageInternalError = "Внутренняя ошибка"

//go:generate mockery --name=telegramClient --dir . --output ./mocks --exported
type telegramClient interface {
	SendMessage(ctx context.Context, userID int64, text string) error
	SendMessageWithoutRemovingKeyboard(ctx context.Context, userID int64, text string) error
	SendKeyboard(ctx context.Context, userID int64, text string, rows [][]string) error
	GetUpdatesChan() <-chan *models.Message
}

//go:generate mockery --name=iterationMessage --dir . --output ./mocks --exported
type iterationMessage interface {
	Iterate(ctx context.Context, message *models.Message, handler MessageHandler, logger log.Logger)
}

type MessageResponse struct {
	Message             string
	Keyboard            [][]string
	DoNotRemoveKeyboard bool
}

type MessageHandler func(ctx context.Context, message *models.Message) (*MessageResponse, error)

type Bot struct {
	tgClient         telegramClient
	iterationMessage iterationMessage
	handlers         map[string]MessageHandler

	logger log.Logger
	cancel context.CancelFunc
	done   chan struct{}
}

func New(tg telegramClient, iterationMessage iterationMessage, logger log.Logger, handlers map[string]MessageHandler) *Bot {
	return &Bot{
		tgClient:         tg,
		iterationMessage: iterationMessage,
		handlers:         handlers,

		logger: logger.With(log.ComponentKey, "Bot"),
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
			handler, ok := b.handlers[message.Text]
			if !ok {
				handler = b.handlers["default"]
			}
			b.iterationMessage.Iterate(ctx, message, handler, b.logger)

		case <-ctx.Done():
			b.logger.WithError(ctx.Err()).Info("bot has been stopped")
			close(b.done)

			return
		}
	}
}

func (b *Bot) UseMiddleware(middleware func(next MessageHandler) MessageHandler) {
	for key, v := range b.handlers {
		b.handlers[key] = middleware(v)
	}
}
