package main

import (
	"context"
	"flag"
	"log"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app/startup"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	exchangeclient "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/exchange"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/telegram"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/repository"
	exchangeservice "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/exchange"
)

func main() {
	configFile := flag.String("config", "", "path to configuration file")
	flag.Parse()

	config, err := startup.NewConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	logger := startup.NewLogger(config.LogLevel)

	tgClient, err := telegram.NewClient(config.Telegram, logger)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to connect to telegram")
	}

	exchangeClient, err := exchangeclient.NewClient(config.ExchangeClient)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to create exchange client")
	}

	dbClient, err := startup.DatabaseConnect(config.Database)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to connect to database")
	}
	defer func() {
		err := dbClient.Close()
		if err != nil {
			logger.WithError(err).
				Warn("failed to close database")
		}
	}()

	userRepo := repository.NewUserRepository(dbClient)
	wasteRepo := repository.NewWasteRepository(dbClient)

	exchangeService, err := exchangeservice.NewService(config.Currency, exchangeClient, logger)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to create exchange repository")
	}

	botComponent := bot.NewBot(tgClient, userRepo, wasteRepo, exchangeService, logger)

	err = app.New(config.App, logger,
		exchangeService,
		botComponent,
	).Run(context.Background())
	if err != nil {
		logger.WithError(err).Fatal("failed during running app")
	}
}
