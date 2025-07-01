package service

import (
	"api/catshelter/internal/domain"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, login, password, name string) (*domain.User, error)
	Login(ctx context.Context, login, password string) (*domain.User, error)
}
