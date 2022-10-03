package startup

import (
	"fmt"
	"os"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/telegram"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram telegram.Config `yaml:"telegram"`

	LogLevel zapcore.Level `yaml:"log_level"`
}

func NewConfig(configFile string) (*Config, error) {
	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading file error: %v", err)
	}

	cfg := &Config{}
	err = yaml.Unmarshal(rawYAML, cfg)
	if err != nil {
		return nil, fmt.Errorf("yaml parsing error: %v", err)
	}

	return cfg, nil
}
