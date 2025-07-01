package service

import "api/catshelter/internal/domain"

type UserService interface {
	FindById(id string) (*domain.User, error)
}
