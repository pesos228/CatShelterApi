package repository

import (
	"api/catshelter/internal/domain"
	"context"
	"errors"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, id string) (*domain.User, error)
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
	FindAll(ctx context.Context) ([]*domain.User, error)
}

var ErrUserNotFound = errors.New("user not found")
