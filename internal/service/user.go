package service

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/repository"
	"context"
	"errors"
	"fmt"
)

type UserService interface {
	FindById(ctx context.Context, id string) (*domain.User, error)
	FindByIdWithCats(ctx context.Context, id string) (*domain.User, error)
	AdoptCat(ctx context.Context, catId, userId string) error
}

type userServiceImpl struct {
	userRepository repository.UserRepository
	catRepository  repository.CatRepository
}

func (u *userServiceImpl) AdoptCat(ctx context.Context, catId, userId string) error {
	cat, err := u.catRepository.FindById(ctx, catId)
	if err != nil {
		if errors.Is(err, repository.ErrCatNotFound) {
			return fmt.Errorf("%w: cat witd id '%s' not found", repository.ErrCatNotFound, catId)
		}
		return err
	}
	_, err = u.FindById(ctx, userId)
	if err != nil {
		return err
	}
	err = cat.AddUser(userId)
	if err != nil {
		return err
	}
	err = u.catRepository.Save(ctx, cat)
	if err != nil {
		return err
	}
	return nil
}

func (u *userServiceImpl) FindById(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepository.FindById(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, fmt.Errorf("%w: user with id '%s' not found", repository.ErrUserNotFound, id)
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

func NewUserService(userRepository repository.UserRepository, catRepository repository.CatRepository) UserService {
	return &userServiceImpl{userRepository: userRepository, catRepository: catRepository}
}
