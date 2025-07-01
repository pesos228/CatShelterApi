package repository

import (
	"api/catshelter/internal/domain"
	"context"
	"errors"
)

type RoleRepository interface {
	Save(ctx context.Context, name string) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

var ErrRoleNotFound = errors.New("role not found")
