package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/api"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type TelegramBot struct {
	config Config
	logger log.Logger

	conn   *grpc.ClientConn
	client api.TelegramBotClient
}

func NewTelegramBot(config Config, logger log.Logger) *TelegramBot {
	return &TelegramBot{
		config: config,
		logger: logger,
	}
}

func (b *TelegramBot) Start() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", b.config.Host, b.config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to grpc on address %s:%d: %w", b.config.Host, b.config.Port, err)
	}

	b.conn = conn
	b.client = api.NewTelegramBotClient(b.conn)

	return nil
}

func (b *TelegramBot) Stop(ctx context.Context) error {
	return b.conn.Close()
}

func (b *TelegramBot) SendMessage(ctx context.Context, userID int64, text string, command enums.CommandType) error {
	_, err := b.client.SendMessage(ctx, &api.Message{
		UserId:  userID,
		Text:    text,
		Command: string(command),
	})
	return err
}
