package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Cat struct {
	BaseModel
	Name   string
	Age    int16
	UserId *string `gorm:"type:uuid"`
}

var ErrValidation = errors.New("validation error")

func NewCat(name string, age int) (*Cat, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("%w: cat must have a name", ErrValidation)
	}
	if age <= 0 {
		return nil, fmt.Errorf("%w: cat age must be positive", ErrValidation)
	}

	return &Cat{
		BaseModel: BaseModel{
			Id: uuid.NewString(),
		},
		Name: name,
		Age:  int16(age),
	}, nil
}

func (c *Cat) AddUser(userId string) error {
	if c.UserId != nil {
		return fmt.Errorf("%w: cat already have a owner", ErrValidation)
	}
	c.UserId = &userId
	return nil
}
