package main

import (
	"context"
	"flag"
	"log"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app/startup"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/bot/handlers"
	exchangeclient "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/exchange"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/telegram"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/http"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/metrics"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/repository"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/cache"
	exchangeservice "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/exchange"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/usercontext"
)

const serviceName = "telegram-bot"

func main() {
	configFile := flag.String("config", "", "path to configuration file")
	flag.Parse()

	config, err := startup.NewConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	logger := startup.NewLogger(serviceName, config.LogLevel)

	tracerProvider, err := metrics.InitTracer(config.Metrics, serviceName)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to create tracer")
	}

	tgClient, err := telegram.NewClient(config.Telegram, logger)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to connect to telegram")
	}

	tgClientDecorator := metrics.NewTelegramClientTracerDecorator(
		metrics.NewTelegramClientLatencyDecorator(tgClient), tracerProvider,
	)

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
		if err := dbClient.Close(); err != nil {
			logger.WithError(err).
				Warn("failed to close database")
		}
	}()

	redisClient, err := startup.RedisConnect(config.Redis)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to connect to redis")
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.WithError(err).
				Warn("failed to close redis")
		}
	}()

	userRepo := metrics.NewUserRepositoryTracerDecorator(
		metrics.NewUserRepositoryAmountErrorsDecorator(
			metrics.NewUserRepositoryLatencyDecorator(
				repository.NewUserRepository(dbClient),
			),
		), tracerProvider,
	)
	wasteRepo := metrics.NewWasteRepositoryTracerDecorator(
		metrics.NewWasteRepositoryAmountErrorsDecorator(
			metrics.NewWasteRepositoryLatencyDecorator(
				repository.NewWasteRepository(dbClient),
			),
		), tracerProvider,
	)

	exchangeService, err := exchangeservice.NewService(config.Currency, exchangeClient, logger)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to create exchange repository")
	}

	userContextService := metrics.NewUserContextServiceTracerDecorator(
		metrics.NewUserContextServiceAmountErrorsDecorator(
			metrics.NewUserContextServiceLatencyDecorator(
				usercontext.NewService(redisClient, config.Currency.Default),
			),
		), tracerProvider,
	)

	cacheService := metrics.NewCacheServiceTracerDecorator(
		metrics.NewCacheServiceAmountErrorsDecorator(
			metrics.NewCacheServiceLatencyDecorator(
				cache.NewService(redisClient, config.Cache),
			),
		), tracerProvider,
	)

	handlers := handlers.NewMessageHandlers(
		userRepo,
		wasteRepo,
		exchangeService,
		userContextService,
	)

	commands := []string{"add", "setLimit", "getLimit", "week", "month", "year", "currency"}

	iterationMessage := metrics.NewIterationMessageTracerDecorator(bot.NewIterationMessage(tgClientDecorator), tracerProvider)
	botComponent := bot.New(tgClientDecorator, iterationMessage, logger, handlers.GetHandlers())
	botComponent.UseMiddleware(bot.CheckUserMiddleware(userRepo))
	botComponent.UseMiddleware(bot.CacheMiddleware(cacheService, logger))
	botComponent.UseMiddleware(bot.LoggerMiddleware(logger))
	botComponent.UseMiddleware(metrics.LatencyMetricMiddleware(commands))
	botComponent.UseMiddleware(metrics.AmountMetricMiddleware(commands))
	botComponent.UseMiddleware(metrics.TracingMiddleware(tracerProvider, commands))

	httpRouter := http.NewHttpRouter(config.Http, logger)

	err = app.New(config.App, logger,
		exchangeService,
		botComponent,
		httpRouter,
	).Run(context.Background())
	if err != nil {
		logger.WithError(err).Fatal("failed during running app")
	}
}
