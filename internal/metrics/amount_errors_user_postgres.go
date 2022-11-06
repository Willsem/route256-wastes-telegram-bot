package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

//go:generate mockery --name=userRepository --dir . --output ./mocks --exported
type userRepository interface {
	UserExists(ctx context.Context, id int64) (bool, error)
	AddUser(ctx context.Context, user *models.User) (*models.User, error)

	SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error)
	GetWasteLimit(ctx context.Context, id int64) (*uint64, error)
}

type UserRepositoryAmountErrorsDecorator struct {
	userRepo    userRepository
	countErrors *prometheus.CounterVec
}

func NewUserRepositoryAmountErrorsDecorator(userRepo userRepository) *UserRepositoryAmountErrorsDecorator {
	return &UserRepositoryAmountErrorsDecorator{
		userRepo: userRepo,
		countErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "count_errors_user_repository",
			Help: "Count of errors in UserRepository methods",
		}, []string{"method"}),
	}
}

func (d *UserRepositoryAmountErrorsDecorator) UserExists(ctx context.Context, id int64) (bool, error) {
	res, err := d.userRepo.UserExists(ctx, id)
	if err != nil {
		d.countErrors.WithLabelValues("UserExists").Inc()
	}
	return res, err
}

func (d *UserRepositoryAmountErrorsDecorator) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	res, err := d.userRepo.AddUser(ctx, user)
	if err != nil {
		d.countErrors.WithLabelValues("AddUser").Inc()
	}
	return res, err
}

func (d *UserRepositoryAmountErrorsDecorator) SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error) {
	res, err := d.userRepo.SetWasteLimit(ctx, id, limit)
	if err != nil {
		d.countErrors.WithLabelValues("SetWasteLimit").Inc()
	}
	return res, err
}

func (d *UserRepositoryAmountErrorsDecorator) GetWasteLimit(ctx context.Context, id int64) (*uint64, error) {
	res, err := d.userRepo.GetWasteLimit(ctx, id)
	if err != nil {
		d.countErrors.WithLabelValues("GetWasteLimit").Inc()
	}
	return res, err
}
