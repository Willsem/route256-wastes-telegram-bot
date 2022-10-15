package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type component interface {
	Start() error
	Stop(ctx context.Context) error
}

type Config struct {
	GracefulTimeout time.Duration `yaml:"graceful_timeout"`
}

type App struct {
	components []component
	config     Config
	logger     log.Logger
}

func New(config Config, logger log.Logger, components ...component) *App {
	return &App{
		components: components,
		config:     config,
		logger:     logger.With(log.ComponentKey, "App"),
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("App is starting")
	defer a.logger.Info("App closed")

	ctxNotify, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g := errgroup.Group{}

	for _, c := range a.components {
		go func(c component) {
			g.Go(c.Start)
		}(c)
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to start components: %w", err)
	}

	a.logger.Info("All components have been started")

	<-ctxNotify.Done()
	a.logger.Info("App received stop signal, all components are stopping")

	if err := a.stop(); err != nil {
		return fmt.Errorf("failed to stop components: %w", err)
	}

	return nil
}

func (a *App) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.GracefulTimeout)
	defer cancel()

	g := errgroup.Group{}

	for _, c := range a.components {
		go func(c component) {
			g.Go(func() error {
				return c.Stop(ctx)
			})
		}(c)
	}

	return g.Wait()
}
