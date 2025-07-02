package service

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/repository"
	"context"
	"fmt"
)

type UserService interface {
	FindById(ctx context.Context, id string) (*domain.User, error)
	FindByIdWithCats(ctx context.Context, id string) (*domain.User, error)
}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func (u *userServiceImpl) FindById(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepository.FindById(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, fmt.Errorf("user with id '%s' not found", id)
		}
		return nil, err
	}
	return user, nil
}

func (u *userServiceImpl) FindByIdWithCats(ctx context.Context, id string) (*domain.User, error) {
	userWithCats, err := u.userRepository.FindByIdWithCats(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, fmt.Errorf("user with id '%s' not found", id)
		}
		return nil, err
	}
	return userWithCats, nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImpl{userRepository: userRepository}
}
