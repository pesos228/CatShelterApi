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
	FindByIdWithAll(ctx context.Context, userId string) (*domain.User, error)
	AdoptCat(ctx context.Context, catId, userId string) error
	AddRole(ctx context.Context, userId, roleName string) error
	RemoveRole(ctx context.Context, userId, roleName string) error
}

type userServiceImpl struct {
	userRepository repository.UserRepository
	catRepository  repository.CatRepository
	roleRepository repository.RoleRepository
}

func (u *userServiceImpl) RemoveRole(ctx context.Context, userId string, roleName string) error {
	user, err := u.userRepository.FindByIdWithRoles(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return fmt.Errorf("%w: user with id '%s' not found", repository.ErrUserNotFound, userId)
		}
		return err
	}
	newRole, err := u.roleRepository.FindByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return fmt.Errorf("%w: role with name '%s' not found", repository.ErrRoleNotFound, roleName)
		}
		return err
	}

	err = user.RemoveRole(newRole)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			return fmt.Errorf("%w: user has no role '%s'", err, roleName)
		}
		if errors.Is(err, domain.ErrCannotRemoveLastRole) {
			return err
		}
		return err
	}

	if err := u.userRepository.UpdateWithRoles(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userServiceImpl) AddRole(ctx context.Context, userId, roleName string) error {
	user, err := u.userRepository.FindByIdWithRoles(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return fmt.Errorf("%w: user with id '%s' not found", repository.ErrUserNotFound, userId)
		}
		return err
	}

	newRole, err := u.roleRepository.FindByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return fmt.Errorf("%w: role with name '%s' not found", repository.ErrRoleNotFound, roleName)
		}
		return err
	}

	if err := user.AddRole(newRole); err != nil {
		return fmt.Errorf("%w: user already have role '%s'", err, roleName)
	}
	if err := u.userRepository.UpdateWithRoles(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userServiceImpl) FindByIdWithAll(ctx context.Context, userId string) (*domain.User, error) {
	user, err := u.userRepository.FindByIdWithAll(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("%w: user with id '%s' not found", repository.ErrUserNotFound, userId)
		}
		return nil, err
	}

	return user, nil
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

func NewUserService(userRepository repository.UserRepository, catRepository repository.CatRepository, roleRepository repository.RoleRepository) UserService {
	return &userServiceImpl{userRepository: userRepository, catRepository: catRepository, roleRepository: roleRepository}
}
