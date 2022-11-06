package metrics

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type UserRepositoryTracerDecorator struct {
	userRepo userRepository
	tracer   trace.Tracer
}

func NewUserRepositoryTracerDecorator(userRepo userRepository, tracerProvider *tracesdk.TracerProvider) *UserRepositoryTracerDecorator {
	return &UserRepositoryTracerDecorator{
		userRepo: userRepo,
		tracer:   tracerProvider.Tracer("user-repository"),
	}
}

func (d *UserRepositoryTracerDecorator) UserExists(ctx context.Context, id int64) (bool, error) {
	ctxTrace, span := d.tracer.Start(ctx, "UserExists")
	defer span.End()

	return d.userRepo.UserExists(ctxTrace, id)
}

func (d *UserRepositoryTracerDecorator) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctxTrace, span := d.tracer.Start(ctx, "AddUser")
	defer span.End()

	return d.userRepo.AddUser(ctxTrace, user)
}

func (d *UserRepositoryTracerDecorator) SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error) {
	ctxTrace, span := d.tracer.Start(ctx, "SetWasteLimit")
	defer span.End()

	return d.userRepo.SetWasteLimit(ctxTrace, id, limit)
}

func (d *UserRepositoryTracerDecorator) GetWasteLimit(ctx context.Context, id int64) (*uint64, error) {
	ctxTrace, span := d.tracer.Start(ctx, "GetWasteLimit")
	defer span.End()

	return d.userRepo.GetWasteLimit(ctxTrace, id)
}
