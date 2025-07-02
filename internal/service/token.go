package service

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
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

type tokenServiceImpl struct {
	auth                   *jwtauth.JWTAuth
	refreshTokenRepository repository.RefreshTokenRepository
	userRepository         repository.UserRepository
}

func (s *tokenServiceImpl) DeleteAllRefreshTokens(ctx context.Context, userId string) error {
	err := s.refreshTokenRepository.DeleteByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return fmt.Errorf("%w: user with id '%s' not found", repository.ErrUserNotFound, userId)
		}
		return err
	}
	return nil
}

func (s *tokenServiceImpl) DeleteRefreshToken(ctx context.Context, token string) error {
	err := s.refreshTokenRepository.DeleteByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return fmt.Errorf("%w: refresh token '%s' not found", repository.ErrRefreshTokenNotFound, token)
		}
		return err
	}
	return nil
}

func (s *tokenServiceImpl) UpdateSession(ctx context.Context, refreshToken string) (*SessionTokens, error) {
	token, err := s.findRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token is expired")
	}

	user, err := s.userRepository.FindByIdWithRoles(ctx, token.UserId)
	if err != nil {
		return nil, fmt.Errorf("db error: %s", err.Error())
	}

	sessionTokens, err := s.generateSessionTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	token.Token = sessionTokens.RefreshToken.Token
	err = s.saveRefreshToken(ctx, &TokenDetails{Id: token.Id, Token: token.Token, UserId: token.UserId, ExpiresAt: sessionTokens.RefreshToken.ExpiresAt})
	if err != nil {
		return nil, err
	}

	return sessionTokens, nil
}

func (s *tokenServiceImpl) CreateSession(ctx context.Context, user *domain.User) (*SessionTokens, error) {
	sessionTokens, err := s.generateSessionTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	err = s.saveRefreshToken(ctx, sessionTokens.RefreshToken)
	if err != nil {
		return nil, err
	}

	return sessionTokens, nil
}

func (s *tokenServiceImpl) findRefreshTokenByToken(ctx context.Context, token string) (*repository.RefreshToken, error) {
	refToken, err := s.refreshTokenRepository.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, fmt.Errorf("%w: refresh token '%s' not found", repository.ErrRefreshTokenNotFound, token)
		}
		return nil, err
	}
	return refToken, nil
}

func (s *tokenServiceImpl) saveRefreshToken(ctx context.Context, token *TokenDetails) error {
	refreshToken := &repository.RefreshToken{
		Id:        token.Id,
		Token:     token.Token,
		UserId:    token.UserId,
		ExpiresAt: token.ExpiresAt,
	}

	err := s.refreshTokenRepository.Save(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("DB error: %w", err)
	}

	return nil
}

func NewTokenService(auth *jwtauth.JWTAuth, refreshTokenRepository repository.RefreshTokenRepository, userRepository repository.UserRepository) TokenService {
	return &tokenServiceImpl{auth: auth, refreshTokenRepository: refreshTokenRepository, userRepository: userRepository}
}

func (s *tokenServiceImpl) generateSessionTokens(ctx context.Context, user *domain.User) (*SessionTokens, error) {
	accessToken, err := s.generateAccessToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}
	refreshTolen, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	return &SessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTolen,
	}, nil
}

func (s *tokenServiceImpl) generateAccessToken(ctx context.Context, user *domain.User) (*TokenDetails, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		exp := time.Now().Add(15 * time.Minute)
		roles := make([]string, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
		claims := map[string]interface{}{
			"user_id": user.Id,
			"roles":   roles,
			"exp":     exp.Unix(),
		}

		_, tokenString, err := s.auth.Encode(claims)
		if err != nil {
			return nil, err
		}

		return &TokenDetails{
			Token:     tokenString,
			ExpiresAt: exp,
		}, err
	}
}

func (s *tokenServiceImpl) generateRefreshToken(ctx context.Context, user *domain.User) (*TokenDetails, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &TokenDetails{
			Token:     uuid.NewString(),
			UserId:    user.Id,
			ExpiresAt: time.Now().Add(24 * time.Hour * 30),
		}, nil
	}
}
