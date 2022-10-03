package app

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app/startup"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/telegram"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/repository"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

func Run(ctx context.Context, config *startup.Config, logger log.Logger) error {
	tgClient, err := telegram.NewClient(config.Telegram, logger)
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %v", err)
	}

	wasteRepo := repository.NewWasteRepository()

	botComponent := bot.NewBot(tgClient, wasteRepo, logger)

	logger.Info("running bot")
	botComponent.Run()

	return nil
}
