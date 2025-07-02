package service

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/repository"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, login, password, name string) (*domain.User, error)
	Login(ctx context.Context, login, password string) (*domain.User, error)
}

type authServiceImpl struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func (s *authServiceImpl) Login(ctx context.Context, login, password string) (*domain.User, error) {
	user, err := s.userRepository.FindByLoginWithRoles(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("user with login: '%s' not found", login)
		}
		return nil, fmt.Errorf("db error: %s", err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}

func (s *authServiceImpl) Register(ctx context.Context, login, password, name string) (*domain.User, error) {
	user, err := s.userRepository.FindByLoginWithRoles(ctx, login)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("db error: %s", err.Error())
		}
	}
	if user != nil {
		return user, fmt.Errorf("user with login '%s' already exists", login)
	}

	user, err = domain.NewUser(login, password, name)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("db error: %s", err.Error())
	}

	role, err := s.roleRepository.FindByName(ctx, "user")
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, errors.New("role 'user' not found")
		}
		return nil, fmt.Errorf("db error: %s", err.Error())
	}

	err = user.AddRole(role)
	if err != nil {
		return user, err
	}

	err = s.userRepository.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("db error: %s", err.Error())
	}
	return user, nil
}

func NewAuthService(userRepository repository.UserRepository, roleRepository repository.RoleRepository) AuthService {
	return &authServiceImpl{userRepository: userRepository, roleRepository: roleRepository}
}
