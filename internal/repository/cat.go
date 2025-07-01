package repository

import (
	"api/catshelter/internal/domain"
	"context"
)

type CatRepository interface {
	Save(ctx context.Context, cat *domain.Cat) error
	FindById(ctx context.Context, id string) (*domain.Cat, error)
	FindAll(ctx context.Context) ([]*domain.Cat, error)
}
