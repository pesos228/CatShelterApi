package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Role struct {
	BaseModel
	Name string
}

func NewRole(name string) (*Role, error) {
	if name == "" {
		return nil, errors.New("role name must not be empty")
	}

	return &Role{
		BaseModel: BaseModel{
			Id: uuid.NewString(),
		},
		Name: name,
	}, nil
}
