package repository

import (
	"api/catshelter/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CatRepository interface {
	Save(ctx context.Context, cat *domain.Cat) error
	FindById(ctx context.Context, id string) (*domain.Cat, error)
	FindWithoutUserId(ctx context.Context, page, pageSize int) ([]*domain.Cat, int64, error)
	FindAll(ctx context.Context) ([]*domain.Cat, error)
}

var ErrCatNotFound = errors.New("cat not found")

type catRepositoryImpl struct {
	db *gorm.DB
}

func (c *catRepositoryImpl) FindWithoutUserId(ctx context.Context, page, pageSize int) ([]*domain.Cat, int64, error) {
	var cats []*domain.Cat
	var count int64

	baseQuery := c.db.WithContext(ctx).Model(&domain.Cat{}).Where("user_id IS NULL")

	if err := baseQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.Scopes(PaginationWithParams(page, pageSize)).Find(&cats).Error; err != nil {
		return nil, 0, err
	}

	return cats, count, nil
}

func (c *catRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Cat, error) {
	var cats []*domain.Cat
	result := c.db.WithContext(ctx).Find(&cats)
	if result.Error != nil {
		return nil, result.Error
	}
	return cats, nil
}

func (c *catRepositoryImpl) FindById(ctx context.Context, id string) (*domain.Cat, error) {
	var cat domain.Cat
	result := c.db.WithContext(ctx).First(&cat, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCatNotFound
		}
		return nil, result.Error
	}
	return &cat, nil
}

func (c *catRepositoryImpl) Save(ctx context.Context, cat *domain.Cat) error {
	return c.db.WithContext(ctx).Save(cat).Error
}

func NewCatRepositoryImpl(db *gorm.DB) CatRepository {
	return &catRepositoryImpl{db: db}
}
