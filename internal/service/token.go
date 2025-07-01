package service

import (
	"api/catshelter/internal/domain"
	"context"
	"time"
)

type TokenService interface {
	CreateSession(ctx context.Context, user *domain.User) (*SessionTokens, error)
	UpdateSession(ctx context.Context, refreshToken string) (*SessionTokens, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllRefreshTokens(ctx context.Context, userId string) error
}

type SessionTokens struct {
	AccessToken  *TokenDetails
	RefreshToken *TokenDetails
}

type TokenDetails struct {
	Id        string
	Token     string
	ExpiresAt time.Time
	UserId    string
}
