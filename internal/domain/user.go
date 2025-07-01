package domain

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Login    string
	Password string
	Name     string
	Role     *Role
	Cats     []*Cat
}

func NewUser(login, password, name string) (*User, error) {
	if len(password) < 8 {
		return nil, errors.New("password is too short")
	}
	if len(login) < 6 {
		return nil, errors.New("login is too short")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		BaseModel: BaseModel{
			Id: uuid.NewString(),
		},
		Login:    login,
		Password: string(hashedPassword),
		Name:     name,
	}, nil
}

func (u *User) SetRole(role *Role) {
	u.Role = role
}
