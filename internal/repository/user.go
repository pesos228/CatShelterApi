package repository

import (
	"api/catshelter/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, id string) (*domain.User, error)
	FindByIdWithAll(ctx context.Context, id string) (*domain.User, error)
	FindByIdWithRoles(ctx context.Context, id string) (*domain.User, error)
	FindByIdWithCats(ctx context.Context, id string) (*domain.User, error)
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
	FindByLoginWithRoles(ctx context.Context, login string) (*domain.User, error)
	FindAll(ctx context.Context) ([]*domain.User, error)
	UpdateWithRoles(ctx context.Context, user *domain.User) error
}

var ErrUserNotFound = errors.New("user not found")

type userRepositoryImpl struct {
	db *gorm.DB
}

func (u *userRepositoryImpl) UpdateWithRoles(ctx context.Context, user *domain.User) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		if err := tx.Model(user).Association("Roles").Replace(user.Roles); err != nil {
			return err
		}

		return nil
	})
}

func (u *userRepositoryImpl) FindByIdWithAll(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Preload("Roles").Preload("Cats").First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *userRepositoryImpl) FindByIdWithRoles(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *userRepositoryImpl) FindByLoginWithRoles(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Preload("Roles").First(&user, "login = ?", login)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *userRepositoryImpl) FindByIdWithCats(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Preload("Cats").First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *userRepositoryImpl) FindAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	result := u.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (u *userRepositoryImpl) FindById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepositoryImpl) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).First(&user, "login = ?", login)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *userRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
	return u.db.WithContext(ctx).Save(user).Error
}

func NewUserReposioryImpl(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
