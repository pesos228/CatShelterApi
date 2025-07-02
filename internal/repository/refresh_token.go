package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	Id        string `gorm:"type:uuid;primary_key;"`
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

type refreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func (r *refreshTokenRepositoryImpl) DeleteByToken(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).Delete(&RefreshToken{}, "token = ?", token)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}
	return nil
}

func (r *refreshTokenRepositoryImpl) DeleteByUserId(ctx context.Context, userId string) error {
	result := r.db.WithContext(ctx).Delete(&RefreshToken{}, "user_id = ?", userId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *refreshTokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	result := r.db.WithContext(ctx).First(&refreshToken, "token = ?", token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, result.Error
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepositoryImpl) Save(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

func NewRefreshTokenRepositoryImpl(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{db: db}
}
