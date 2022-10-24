package startup

import (
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app"
	exchangeclient "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/exchange"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/telegram"
	exchangeservice "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/service/exchange"
)

type Config struct {
	App            app.Config             `yaml:"app"`
	Telegram       telegram.Config        `yaml:"telegram"`
	ExchangeClient exchangeclient.Config  `yaml:"exchange_client"`
	Currency       exchangeservice.Config `yaml:"currency"`
	Database       DatabaseConfig         `yaml:"database"`

	LogLevel zapcore.Level `yaml:"log_level"`
}

func NewConfig(configFile string) (*Config, error) {
	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading file error: %w", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(rawYAML, cfg); err != nil {
		return nil, fmt.Errorf("yaml parsing error: %w", err)
	}

	return cfg, nil
}
