package bot

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
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

func TracingMiddleware() MessageMiddleware {
	middleware := func(next MessageHandler) MessageHandler {
		return func(ctx context.Context, message *models.Message) (*MessageResponse, error) {
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
