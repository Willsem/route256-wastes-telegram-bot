package usercontext

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models/enums"
)

// TODO: переписать на Redis в 6 дз
type Service struct {
	users map[int64]enums.UserContext

	defaultCurrency string
	currencies      map[int64]string
}

func NewService(defaultCurrency string) *Service {
	return &Service{
		users: make(map[int64]enums.UserContext),

		defaultCurrency: defaultCurrency,
		currencies:      make(map[int64]string),
	}
}

func (s *Service) SetContext(ctx context.Context, userID int64, context enums.UserContext) error {
	s.users[userID] = context
	return nil
}

func (s *Service) GetContext(ctx context.Context, userID int64) (enums.UserContext, error) {
	if _, ok := s.users[userID]; !ok {
		s.users[userID] = enums.NoContext
		return enums.NoContext, nil
	}

	return s.users[userID], nil
}

func (s *Service) SetCurrency(ctx context.Context, userID int64, currency string) error {
	s.currencies[userID] = currency
	return nil
}

func (s *Service) GetCurrency(ctx context.Context, userID int64) (string, error) {
	if _, ok := s.currencies[userID]; !ok {
		s.currencies[userID] = s.defaultCurrency
		return s.defaultCurrency, nil
	}

	return s.currencies[userID], nil
}
