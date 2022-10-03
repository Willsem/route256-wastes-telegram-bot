package zap

import (
	"os"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const errKey = "error"

type LoggerImpl struct {
	logger *zap.SugaredLogger
}

func NewLogger(loggerName string, logLevel zapcore.Level) log.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	logger := zap.New(core,
		zap.Fields(zap.String("name", loggerName)),
	).Sugar()

	return &LoggerImpl{
		logger: logger,
	}
}

func (log *LoggerImpl) With(args ...interface{}) log.Logger {
	return &LoggerImpl{
		logger: log.logger.With(args...),
	}
}

func (log *LoggerImpl) WithError(err error) log.Logger {
	return log.With(errKey, err)
}

func (log *LoggerImpl) Debug(args ...interface{}) {
	log.logger.Debug(args...)
}

func (log *LoggerImpl) Info(args ...interface{}) {
	log.logger.Info(args...)
}

func (log *LoggerImpl) Warn(args ...interface{}) {
	log.logger.Warn(args...)
}

func (log *LoggerImpl) Error(args ...interface{}) {
	log.logger.Error(args...)
}

func (log *LoggerImpl) Fatal(args ...interface{}) {
	log.logger.Fatal(args...)
}

func (log *LoggerImpl) Debugf(template string, args ...interface{}) {
	log.logger.Debugf(template, args...)
}

func (log *LoggerImpl) Infof(template string, args ...interface{}) {
	log.logger.Infof(template, args...)
}

func (log *LoggerImpl) Warnf(template string, args ...interface{}) {
	log.logger.Warnf(template, args...)
}

func (log *LoggerImpl) Errorf(template string, args ...interface{}) {
	log.logger.Errorf(template, args...)
}

func (log *LoggerImpl) Fatalf(template string, args ...interface{}) {
	log.logger.Fatalf(template, args...)
}
