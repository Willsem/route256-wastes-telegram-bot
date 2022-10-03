package main

import (
	"context"
	"flag"
	"log"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app/startup"
)

func main() {
	configFile := flag.String("config", "", "path to configuration file")
	flag.Parse()

	config, err := startup.NewConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	logger := startup.NewLogger(config.LogLevel)

	ctx := context.Background()
	if err = app.Run(ctx, config, logger); err != nil {
		logger.
			WithError(err).
			Fatal("failed to run application")
	}
}
