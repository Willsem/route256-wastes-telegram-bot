package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

type UserRepositoryLatencyDecorator struct {
	userRepo userRepository
	latency  *prometheus.HistogramVec
}

func NewUserRepositoryLatencyDecorator(userRepo userRepository) *UserRepositoryLatencyDecorator {
	return &UserRepositoryLatencyDecorator{
		userRepo: userRepo,
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "latency_user_repository",
			Help:    "Duration of UserRepository methods",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0},
		}, []string{"method"}),
	}
}

func (d *UserRepositoryLatencyDecorator) UserExists(ctx context.Context, id int64) (bool, error) {
	startTime := time.Now()
	res, err := d.userRepo.UserExists(ctx, id)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("UserExists").Observe(duration.Seconds())

	return res, err
}

func (d *UserRepositoryLatencyDecorator) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	startTime := time.Now()
	res, err := d.userRepo.AddUser(ctx, user)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("AddUser").Observe(duration.Seconds())

	return res, err
}

func (d *UserRepositoryLatencyDecorator) SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error) {
	startTime := time.Now()
	res, err := d.userRepo.SetWasteLimit(ctx, id, limit)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("SetWasteLimit").Observe(duration.Seconds())

	return res, err
}

func (d *UserRepositoryLatencyDecorator) GetWasteLimit(ctx context.Context, id int64) (*uint64, error) {
	startTime := time.Now()
	res, err := d.userRepo.GetWasteLimit(ctx, id)
	duration := time.Since(startTime)

	d.latency.WithLabelValues("GetWasteLimit").Observe(duration.Seconds())

	return res, err
}
