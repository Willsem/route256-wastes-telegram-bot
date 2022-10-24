package repository

import (
	"context"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/user"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	model, err := r.client.User.
		Create().
		SetID(user.ID).
		SetFirstName(user.FirstName).
		SetLastName(user.LastName).
		SetUserName(user.UserName).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &models.User{
		User: model,
	}, nil
}

func (r *UserRepository) UserExists(ctx context.Context, id int64) (bool, error) {
	exists, err := r.client.User.Query().
		Where(user.ID(id)).
		Exist(ctx)
	if err != nil {
		return false, err
	}

	return exists, err
}

func (r *UserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	model, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		User: model,
	}, nil
}

func (r *UserRepository) SetWasteLimit(ctx context.Context, id int64, limit uint64) (*models.User, error) {
	model, err := r.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	updated, err := model.Update().SetWasteLimit(limit).Save(ctx)
	if err != nil {
		return nil, err
	}

	return &models.User{
		User: updated,
	}, nil
}

func (r *UserRepository) GetWasteLimit(ctx context.Context, id int64) (*uint64, error) {
	model, err := r.client.User.Query().
		Select(user.FieldWasteLimit).
		Where(user.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return model.WasteLimit, nil
}
