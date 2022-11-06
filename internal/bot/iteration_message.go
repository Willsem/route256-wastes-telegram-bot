package bot

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type IterationMessage struct {
	tgClient telegramClient
}

func NewIterationMessage(tgClient telegramClient) *IterationMessage {
	return &IterationMessage{
		tgClient: tgClient,
	}
}

func (i *IterationMessage) Iterate(ctx context.Context, message *models.Message, handler MessageHandler, logger log.Logger) {
	response, err := handler(ctx, message)
	if err != nil {
		logger.WithError(err).
			With("message", message).
			Error("failed to respond the message")
		err := i.tgClient.SendMessage(ctx, message.From.ID, messageInternalError)
		if err != nil {
			logger.WithError(err).
				With("response", response).
				With("message", message).
				Error("failed to send the message")
		}
		return
	}

	if response.Keyboard != nil {
		err = i.tgClient.SendKeyboard(ctx, message.From.ID, response.Message, response.Keyboard)
	} else if response.DoNotRemoveKeyboard {
		err = i.tgClient.SendMessageWithoutRemovingKeyboard(ctx, message.From.ID, response.Message)
	} else {
		err = i.tgClient.SendMessage(ctx, message.From.ID, response.Message)
	}

	if err != nil {
		logger.WithError(err).
			With("response", response).
			With("message", message).
			Error("failed to send the message")
	}
}
