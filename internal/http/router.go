package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type Config struct {
	Port int `yaml:"port"`
}

type HttpRouter struct {
	port   int
	server *http.Server
	logger log.Logger
}

func NewHttpRouter(config Config, logger log.Logger) *HttpRouter {
	serveMux := http.NewServeMux()
	serveMux.Handle("/metrics", promhttp.Handler())

	return &HttpRouter{
		port: config.Port,
		server: &http.Server{
			Handler: serveMux,
		},
		logger: logger.With(log.ComponentKey, "Http server"),
	}
}

func (r *HttpRouter) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		return fmt.Errorf("failed to listen the port %d: %w", r.port, err)
	}

	go func() {
		r.logger.Infof("server is listening the port %d", r.port)

		if err := r.server.Serve(listener); err != http.ErrServerClosed {
			r.logger.WithError(err).
				Fatalf("fail to serve the server on the port %d", r.port)
		}
	}()
	return nil
}

func (r *HttpRouter) Stop(ctx context.Context) error {
	r.logger.Info("http router is stopping")
	return r.server.Shutdown(ctx)
}
