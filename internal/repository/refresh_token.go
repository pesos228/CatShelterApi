package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	Id        string
	Token     string
	UserId    string
	ExpiresAt time.Time
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserId(ctx context.Context, userId string) error
}

func (s *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Id == "" {
		s.Id = uuid.NewString()
	}
	return
}

var ErrRefreshTokenNotFound = errors.New("refresh token not found")
