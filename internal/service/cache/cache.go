package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

type Config struct {
	Expiration time.Duration `yaml:"expiration"`
}

type Service struct {
	client *redis.Client
	config Config
}

func NewService(client *redis.Client, config Config) *Service {
	return &Service{
		client: client,
		config: config,
	}
}

func getKey(userID int64, command enums.CommandType) string {
	return fmt.Sprintf("cache_%d_%s", userID, command)
}

func (s *Service) Set(ctx context.Context, userID int64, command enums.CommandType, value string) error {
	err := s.client.Set(ctx, getKey(userID, command), value, s.config.Expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value to the cache: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, userID int64, command enums.CommandType) (string, error) {
	value, err := s.client.Get(ctx, getKey(userID, command)).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get value from the cache: %w", err)
	}

	return value, nil
}

func (s *Service) Clear(ctx context.Context, userID int64, command enums.CommandType) error {
	err := s.client.Del(ctx, getKey(userID, command)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete a key from the cache: %w", err)
	}

	return nil
}

func (s *Service) ClearKeys(ctx context.Context, userID int64, commands ...enums.CommandType) error {
	for _, command := range commands {
		err := s.Clear(ctx, userID, command)
		if err != nil {
			return err
		}
	}

	return nil
}
