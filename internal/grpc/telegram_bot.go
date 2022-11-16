package grpc

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/api"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

//go:generate mockery --name=telegramClient --dir . --output ./mocks --exported
type telegramClient interface {
	SendMessage(ctx context.Context, userID int64, text string) error
}

type cacheService interface {
	Set(ctx context.Context, userID int64, command enums.CommandType, value string) error
}

type TelegramBotClient struct {
	api.UnimplementedTelegramBotServer

	tgClient telegramClient
	cache    cacheService
}

func NewTelegramBotClient(tgClient telegramClient, cacheService cacheService) *TelegramBotClient {
	return &TelegramBotClient{
		tgClient: tgClient,
		cache:    cacheService,
	}
}

func (c *TelegramBotClient) SendMessage(ctx context.Context, msg *api.Message) (*api.EmptyMessage, error) {
	err := c.tgClient.SendMessage(ctx, msg.GetUserId(), msg.GetText())
	if err != nil {
		return nil, fmt.Errorf("failed to send message by tg client: %w", err)
	}

	err = c.cache.Set(ctx, msg.GetUserId(), enums.CommandType(msg.GetCommand()), msg.GetText())
	if err != nil {
		return nil, fmt.Errorf("failed to set value to the cache: %w", err)
	}

	return &api.EmptyMessage{}, nil
}
