package startup

import (
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(loggerName string, logLevel zapcore.Level) log.Logger {
	return zap.NewLogger(loggerName, logLevel)
}
