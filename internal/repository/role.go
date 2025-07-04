package repository

import (
	"api/catshelter/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(ctx context.Context, role *domain.Role) error
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

var ErrRoleNotFound = errors.New("role not found")

type roleRepositoryImpl struct {
	db *gorm.DB
}

func (r *roleRepositoryImpl) Save(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepositoryImpl) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	result := r.db.WithContext(ctx).First(&role, "LOWER(name) = LOWER(?)", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrRoleNotFound
	}
	return &role, nil
}

func NewRoleRepositoryImpl(db *gorm.DB) RoleRepository {
	return &roleRepositoryImpl{db: db}
}
