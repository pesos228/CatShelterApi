package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Login    string `gorm:"unique"`
	Password string
	Name     string
	Roles    []*Role `gorm:"many2many:user_roles;"`
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

func (u *User) AddRole(role *Role) error {
	_, ok := u.IsHaveRole(role)
	if ok {
		return fmt.Errorf("user already have role '%s'", role.Name)
	}
	u.Roles = append(u.Roles, role)
	return nil
}

func (u *User) RemoveRole(role *Role) error {
	index, ok := u.IsHaveRole(role)
	if !ok {
		return fmt.Errorf("user has no role '%s'", role.Name)
	}
	u.Roles[index] = u.Roles[len(u.Roles)-1]
	u.Roles = u.Roles[:len(u.Roles)-1]
	return nil
}

func (u *User) IsHaveRole(role *Role) (int, bool) {
	for i, r := range u.Roles {
		if strings.EqualFold(r.Name, role.Name) {
			return i, true
		}
	}
	return -1, false
}
