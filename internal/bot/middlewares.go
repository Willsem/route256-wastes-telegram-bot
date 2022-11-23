package bot

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type MessageMiddleware func(next MessageHandler) MessageHandler

func LoggerMiddleware(l log.Logger) MessageMiddleware {
	logger := l.With(log.ComponentKey, "Logger middleware")

	middleware := func(next MessageHandler) MessageHandler {
		return func(ctx context.Context, message *models.Message) (*MessageResponse, error) {
			logger.Infof("Message \"%s\" from user %s(%d)", message.Text, message.From.UserName, message.From.ID)
			return next(ctx, message)
		}
	}

	return middleware
}

//go:generate mockery --name=userRepository --dir . --output ./mocks --exported
type userRepository interface {
	UserExists(ctx context.Context, id int64) (bool, error)
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
}

// CheckUserMiddleware middleware for adding new users
// and make sure that user exists during the running message handler.
func CheckUserMiddleware(userRepo userRepository) MessageMiddleware {
	middleware := func(next MessageHandler) MessageHandler {
		return func(ctx context.Context, message *models.Message) (*MessageResponse, error) {
			exists, err := userRepo.UserExists(ctx, message.From.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check exising user: %w", err)
			}

			if !exists {
				_, err := userRepo.AddUser(ctx, message.From)
				if err != nil {
					return nil, fmt.Errorf("failed to adding user: %w", err)
				}
			}

			return next(ctx, message)
		}
	}

	return middleware
}

//go:generate mockery --name=cacheService --dir . --output ./mocks --exported
type cacheService interface {
	Set(ctx context.Context, userID int64, command enums.CommandType, value string) error
	Get(ctx context.Context, userID int64, command enums.CommandType) (string, error)
	Clear(ctx context.Context, userID int64, command enums.CommandType) error
	ClearKeys(ctx context.Context, userID int64, commands ...enums.CommandType) error
}

func CacheMiddleware(cacheService cacheService, logger log.Logger) MessageMiddleware {
	logger = logger.With(log.ComponentKey, "Cache middleware")
	middleware := func(next MessageHandler) MessageHandler {
		return func(ctx context.Context, message *models.Message) (*MessageResponse, error) {
			command, err := enums.ParseCommandType(message.Text)
			if err != nil {
				return next(ctx, message)
			}

			switch command {
			case enums.CommandTypeCurrency:
				err = cacheService.Clear(ctx, message.From.ID, enums.CommandTypeGetLimit)
				if err != nil {
					logger.WithError(err).
						Info("failed to clear key in the cache")
				}
				fallthrough

			case enums.CommandTypeAdd:
				err = cacheService.ClearKeys(ctx, message.From.ID,
					enums.CommandTypeWeekReport,
					enums.CommandTypeMonthReport,
					enums.CommandTypeYearReport,
				)
				if err != nil {
					logger.WithError(err).
						Info("failed to clear key in the cache")
				}
				return next(ctx, message)

			case enums.CommandTypeSetLimit:
				err = cacheService.Clear(ctx, message.From.ID, enums.CommandTypeGetLimit)
				if err != nil {
					logger.WithError(err).
						Info("failed to clear key in the cache")
				}
				return next(ctx, message)

			default:
				result, err := cacheService.Get(ctx, message.From.ID, command)
				if err == nil {
					return &MessageResponse{
						Message: result,
					}, nil
				}
			}

			resp, err := next(ctx, message)

			if command != enums.CommandTypeWeekReport &&
				command != enums.CommandTypeMonthReport &&
				command != enums.CommandTypeYearReport {
				if err := cacheService.Set(ctx, message.From.ID, command, resp.Message); err != nil {
					logger.WithError(err).
						Error("failed to upload the message to the cache")
				}
			}

			return resp, err
		}
	}

	return middleware
}
