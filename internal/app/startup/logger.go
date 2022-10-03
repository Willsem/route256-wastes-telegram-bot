package startup

import (
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log/zap"
	"go.uber.org/zap/zapcore"
)

const loggerName = "telegram-bot"

func NewLogger(logLevel zapcore.Level) log.Logger {
	return zap.NewLogger(loggerName, logLevel)
}
