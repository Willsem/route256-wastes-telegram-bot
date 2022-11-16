package startup

import (
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/grpc"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/http"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/metrics"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/kafka"
)

type ReportServiceConfig struct {
	App      app.Config           `yaml:"app"`
	Database DatabaseConfig       `yaml:"database"`
	Kafka    KafkaConfig          `yaml:"kafka"`
	Consumer kafka.ConsumerConfig `yaml:"consumer"`
	Http     http.Config          `yaml:"http"`
	Grpc     grpc.Config          `yaml:"grpc_client"`
	Metrics  metrics.Config       `yaml:"metrics"`

	LogLevel zapcore.Level `yaml:"log_level"`
}

func NewReportServiceConfig(configFile string) (*ReportServiceConfig, error) {
	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading file error: %w", err)
	}

	cfg := &ReportServiceConfig{}
	if err = yaml.Unmarshal(rawYAML, cfg); err != nil {
		return nil, fmt.Errorf("yaml parsing error: %w", err)
	}

	return cfg, nil
}
