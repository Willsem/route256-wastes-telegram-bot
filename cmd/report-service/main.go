package main

import (
	"context"
	"flag"
	"log"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app/startup"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/grpc"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/http"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/metrics"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/repository"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/kafka"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/wastereport"
)

const serviceName = "report-service"

func main() {
	configFile := flag.String("config", "", "path to configuration file")
	flag.Parse()

	config, err := startup.NewReportServiceConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	logger := startup.NewLogger(serviceName, config.LogLevel)

	logger.With("config", config).Info("application staring with this config")

	tracerProvider, err := metrics.InitTracer(config.Metrics, serviceName)
	if err != nil {
		logger.WithError(err).
			Fatal("failed to create tracer")
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

	kafkaClient := startup.NewKafkaConsumer(config.Kafka)
	defer func() {
		if err := kafkaClient.Close(); err != nil {
			logger.WithError(err).
				Warn("failed to close kafka connection")
		}
	}()

	wasteRepo := metrics.NewWasteRepositoryTracerDecorator(
		metrics.NewWasteRepositoryAmountErrorsDecorator(
			metrics.NewWasteRepositoryLatencyDecorator(
				repository.NewWasteRepository(dbClient),
			),
		), tracerProvider,
	)

	consumerComponent := kafka.NewConsumer(kafkaClient, config.Consumer, logger)

	httpRouter := http.NewHttpRouter(config.Http, logger)
	grpcClient := grpc.NewTelegramBot(config.Grpc, logger)

	reportService := wastereport.NewService(consumerComponent, wasteRepo, grpcClient, logger)

	err = app.New(config.App, logger,
		consumerComponent,
		httpRouter,
		grpcClient,
		reportService,
	).Run(context.Background())
	if err != nil {
		logger.WithError(err).Fatal("failed during running app")
	}
}
