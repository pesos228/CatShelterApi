package service

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/repository"
	"context"
	"fmt"
)

type CatService interface {
	FindLonelyCats(ctx context.Context) ([]*domain.Cat, error)
	AddCat(ctx context.Context, name string, age int) error
}

type catServiceImpl struct {
	catRepository repository.CatRepository
}

func (c *catServiceImpl) AddCat(ctx context.Context, name string, age int) error {
	newCat, err := domain.NewCat(name, age)
	if err != nil {
		return err
	}
	err = c.catRepository.Save(ctx, newCat)
	if err != nil {
		return fmt.Errorf("DB error: %s", err.Error())
	}

	return nil
}

func (c *catServiceImpl) FindLonelyCats(ctx context.Context) ([]*domain.Cat, error) {
	lonelyCats, err := c.catRepository.FindWithoutUserId(ctx)
	if err != nil {
		return nil, fmt.Errorf("DB error: %s", err.Error())
	}
	return lonelyCats, nil
}

func NewCatService(catRepository repository.CatRepository) CatService {
	return &catServiceImpl{catRepository: catRepository}
}
