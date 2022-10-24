package exchange

import (
	"context"
	"errors"
	"sync"
	"time"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type Config struct {
	Default            string `yaml:"default"`
	DesignationDefault string `yaml:"designation_default"`

	Used            []string `yaml:"used"`
	DesignationUsed []string `yaml:"designation_used"`

	UpdateTimeout time.Duration `yaml:"update_timeout"`
	RetryTimeout  time.Duration `yaml:"retry_timeout"`
}

//go:generate mockery --name=exchangeClient --dir . --output ./mocks --exported
type exchangeClient interface {
	GetExchange(ctx context.Context, base string, symbols []string) (*models.ExchangeData, error)
}

var (
	ErrIncorrectConfig  = errors.New("different length of designations and currencies in a configs")
	ErrCurrencyNotFound = errors.New("currency not found in repository")
	ErrDataNotPrepared  = errors.New("currency exchange does not prepared")
)

type Service struct {
	exchangeClient exchangeClient
	config         Config
	logger         log.Logger

	data         map[string]float64
	designations map[string]string
	mutex        *sync.RWMutex

	cancel context.CancelFunc
	done   chan struct{}
}

func NewService(config Config, exchangeClient exchangeClient, logger log.Logger) (*Service, error) {
	if len(config.Used) != len(config.DesignationUsed) {
		return nil, ErrIncorrectConfig
	}

	s := &Service{
		exchangeClient: exchangeClient,
		config:         config,
		logger:         logger.With(log.ComponentKey, "Exchange service"),
		mutex:          &sync.RWMutex{},
	}

	designations := make(map[string]string)
	designations[config.Default] = config.DesignationDefault
	for i := range config.Used {
		designations[config.Used[i]] = config.DesignationUsed[i]
	}

	s.designations = designations

	return s, nil
}

func (s *Service) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	s.cancel = cancel
	s.done = make(chan struct{})

	go s.updateDataByTicker(ctx)

	return s.updateData(ctx)
}

func (s *Service) Stop(ctx context.Context) error {
	s.cancel()

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Service) GetDefaultCurrency() string {
	return s.config.Default
}

func (s *Service) GetUsedCurrencies() []string {
	return s.config.Used
}

func (s *Service) GetDesignation(currency string) (string, error) {
	s.mutex.RLock()
	data := s.designations
	s.mutex.RUnlock()

	designation, ok := data[currency]
	if !ok {
		return "", ErrCurrencyNotFound
	}

	return designation, nil
}

func (s *Service) GetExchange(currency string) (float64, error) {
	s.mutex.RLock()
	data := s.data
	s.mutex.RUnlock()

	if data == nil {
		return 0, ErrDataNotPrepared
	}

	exchange, ok := data[currency]
	if !ok {
		return 0, ErrCurrencyNotFound
	}

	return exchange, nil
}

func (s *Service) updateDataByTicker(ctx context.Context) {
	ticker := time.NewTicker(s.config.UpdateTimeout)

	for {
		select {
		case <-ticker.C:
			err := s.updateData(ctx)
			if err != nil {
				s.logger.WithError(err).Warn("failed to update exchanges data")
				ticker.Reset(s.config.RetryTimeout)
			} else {
				ticker.Reset(s.config.UpdateTimeout)
			}

		case <-ctx.Done():
			s.logger.WithError(ctx.Err()).Info("exchange repository has been closed")
			close(s.done)

			return
		}
	}
}

func (s *Service) updateData(ctx context.Context) error {
	s.logger.Info("try to update exchange data")

	data, err := s.exchangeClient.GetExchange(ctx, s.config.Default, s.config.Used)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	if s.data == nil {
		s.data = make(map[string]float64, len(data.Rates))
	}
	for key, value := range data.Rates {
		s.data[key] = value
	}
	s.data[s.config.Default] = 1.0
	s.mutex.Unlock()

	s.logger.
		With("exchange data", data).
		Info("exchange data updated successfully")

	return nil
}
