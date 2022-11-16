package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/api"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type Config struct {
	Port int `yaml:"port"`
}

type Server struct {
	port   int
	logger log.Logger

	server *grpc.Server

	tgClient telegramClient
	cache    cacheService

	isClosed bool
}

func NewServer(config Config, tgClient telegramClient, cacheService cacheService, logger log.Logger) *Server {
	return &Server{
		port:   config.Port,
		logger: logger.With(log.ComponentKey, "Grpc server"),

		tgClient: tgClient,
		cache:    cacheService,

		isClosed: false,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen port %d: %w", s.port, err)
	}

	s.server = grpc.NewServer()
	api.RegisterTelegramBotServer(s.server, NewTelegramBotClient(s.tgClient, s.cache))

	go func() {
		s.logger.Infof("server is listening the port %d", s.port)

		if err := s.server.Serve(listener); err != nil && !s.isClosed {
			s.logger.WithError(err).
				Fatalf("fail to serve the server on the port %d", s.port)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("grpc server is stopping")
	s.isClosed = true
	s.server.Stop()
	return nil
}
