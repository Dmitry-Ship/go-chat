package services

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type authService struct {
	users    domain.UserRepository
	jwTokens JWTokens
}

func NewAuthService(users domain.UserRepository, config config.Auth) *authService {
	return &authService{
		users:    users,
		jwTokens: NewJWTokens(config),
	}
}

func (a *authService) Login(ctx context.Context, username string, password string) (Tokens, error) {
	user, err := a.users.FindByUsername(ctx, username)

	if err != nil {
		return Tokens{}, fmt.Errorf("find by username error: %w", err)
	}

	if err := domain.ComparePassword(user.PasswordHash, password); err != nil {
		return Tokens{}, fmt.Errorf("compare password error: %w", err)
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return Tokens{}, fmt.Errorf("create tokens error: %w", err)
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Update(ctx, user); err != nil {
		return Tokens{}, fmt.Errorf("update user error: %w", err)
	}

	return newTokens, err
}

func (a *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	user, err := a.users.GetByID(ctx, userID)

	if err != nil {
		return fmt.Errorf("get by id error: %w", err)
	}

	user.SetRefreshToken("")

	err = a.users.Update(ctx, user)

	if err != nil {
		return fmt.Errorf("update user error: %w", err)
	}

	return nil
}

func (a *authService) SignUp(ctx context.Context, username string, password string) (Tokens, error) {
	if err := domain.ValidateUsername(username); err != nil {
		return Tokens{}, fmt.Errorf("validate username error: %w", err)
	}

	hashedPassword, err := domain.HashPassword(password)

	if err != nil {
		return Tokens{}, fmt.Errorf("hash password error: %w", err)
	}

	userID := uuid.New()

	user := domain.NewUser(userID, username, hashedPassword)

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return Tokens{}, fmt.Errorf("create tokens error: %w", err)
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Store(ctx, user); err != nil {
		return Tokens{}, fmt.Errorf("store user error: %w", err)
	}

	return newTokens, nil
}

func (a *authService) RotateTokens(ctx context.Context, refreshTokenString string) (Tokens, error) {
	userID, err := a.jwTokens.ParseRefreshToken(refreshTokenString)

	if err != nil {
		return Tokens{}, fmt.Errorf("parse refresh token error: %w", err)
	}

	user, err := a.users.GetByID(ctx, userID)

	if err != nil {
		return Tokens{}, fmt.Errorf("get by id error: %w", err)
	}

	if user.RefreshToken != refreshTokenString {
		return Tokens{}, fmt.Errorf("invalid token")
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return Tokens{}, fmt.Errorf("create tokens error: %w", err)
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Update(ctx, user); err != nil {
		return Tokens{}, fmt.Errorf("update user error: %w", err)
	}

	return newTokens, nil
}

func (a *authService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	return a.jwTokens.ParseAccessToken(tokenString)
}
