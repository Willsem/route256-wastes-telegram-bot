package usercontext

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

const (
	userContext  = "usercontext"
	userCurrency = "usercurrency"
)

type Service struct {
	client *redis.Client

	defaultCurrency string
	currencies      map[int64]string
}

func getKey(userID int64, key string) string {
	return fmt.Sprintf("user_context_%s_%d", key, userID)
}

func NewService(client *redis.Client, defaultCurrency string) *Service {
	return &Service{
		client: client,

		defaultCurrency: defaultCurrency,
		currencies:      make(map[int64]string),
	}
}

func (s *Service) SetContext(ctx context.Context, userID int64, context enums.UserContext) error {
	err := s.client.Set(ctx, getKey(userID, userContext), int(context), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set user context: %w", err)
	}

	return nil
}

func (s *Service) GetContext(ctx context.Context, userID int64) (enums.UserContext, error) {
	userContext, err := s.client.Get(ctx, getKey(userID, userContext)).Int()
	if err == redis.Nil {
		err := s.SetContext(ctx, userID, enums.NoContext)
		if err != nil {
			return enums.NoContext, fmt.Errorf("failed to set context after receiving nil of getting: %w", err)
		}

		return enums.NoContext, nil
	}
	if err != nil {
		return enums.NoContext, fmt.Errorf("failed to set user context: %w", err)
	}

	return enums.UserContext(userContext), nil
}

func (s *Service) SetCurrency(ctx context.Context, userID int64, currency string) error {
	err := s.client.Set(ctx, getKey(userID, userCurrency), currency, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set user currency: %w", err)
	}

	return nil
}

func (s *Service) GetCurrency(ctx context.Context, userID int64) (string, error) {
	currency, err := s.client.Get(ctx, getKey(userID, userCurrency)).Result()
	if err == redis.Nil {
		err := s.SetCurrency(ctx, userID, s.defaultCurrency)
		if err != nil {
			return "", fmt.Errorf("failed to set currency after receiving nil of getting: %w", err)
		}

		return s.defaultCurrency, nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user currency: %w", err)
	}

	return currency, nil
}
