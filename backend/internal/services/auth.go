package services

import (
	"fmt"

	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	users    domain.UserRepository
	jwTokens JWTokens
}

type AuthService interface {
	Login(username string, password string) (tokens, error)
	Logout(userID uuid.UUID) error
	SignUp(username string, password string) (tokens, error)
	RotateTokens(refreshTokenString string) (tokens, error)
	ParseAccessToken(accessTokenString string) (uuid.UUID, error)
}

func NewAuthService(users domain.UserRepository, config config.Auth) *authService {
	return &authService{
		users:    users,
		jwTokens: NewJWTokens(config),
	}
}

func (a *authService) Login(username string, password string) (tokens, error) {
	user, err := a.users.FindByUsername(username)

	if err != nil {
		return tokens{}, fmt.Errorf("find by username error: %w", err)
	}

	userPassword, err := domain.NewUserPassword(password, func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	if err != nil {
		return tokens{}, fmt.Errorf("new user password error: %w", err)
	}

	if err = user.Password.Compare(userPassword, func(p1 []byte, p2 []byte) error {
		return bcrypt.CompareHashAndPassword(p1, p2)
	}); err != nil {
		return tokens{}, fmt.Errorf("compare password error: %w", err)
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, fmt.Errorf("create tokens error: %w", err)
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Update(user); err != nil {
		return tokens{}, fmt.Errorf("update user error: %w", err)
	}

	return newTokens, err
}

func (a *authService) Logout(userID uuid.UUID) error {
	user, err := a.users.GetByID(userID)

	if err != nil {
		return fmt.Errorf("get by id error: %w", err)
	}

	user.SetRefreshToken("")

	err = a.users.Update(user)

	if err != nil {
		return fmt.Errorf("update user error: %w", err)
	}

	return nil
}

func (a *authService) SignUp(username string, password string) (tokens, error) {
	name, err := domain.NewUserName(username)

	if err != nil {
		return tokens{}, fmt.Errorf("new user name error: %w", err)
	}

	hashedPassword, err := domain.NewUserPassword(password, func(p []byte) ([]byte, error) {
		return bcrypt.GenerateFromPassword(p, 14)
	})

	if err != nil {
		return tokens{}, fmt.Errorf("new user password error: %w", err)
	}

	userID := uuid.New()

	user := domain.NewUser(userID, name, hashedPassword)

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, fmt.Errorf("create tokens error: %w", err)
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Store(user); err != nil {
		return tokens{}, fmt.Errorf("store user error: %w", err)
	}

	return newTokens, nil
}

func (a *authService) RotateTokens(refreshTokenString string) (tokens, error) {
	userID, err := a.jwTokens.ParseRefreshToken(refreshTokenString)

	if err != nil {
		return tokens{}, fmt.Errorf("parse refresh token error: %w", err)
	}

	user, err := a.users.GetByID(userID)

	if err != nil {
		return tokens{}, fmt.Errorf("get by id error: %w", err)
	}

	if user.RefreshToken != refreshTokenString {
		return tokens{}, fmt.Errorf("invalid token")
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, fmt.Errorf("create tokens error: %w", err)
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Update(user); err != nil {
		return tokens{}, fmt.Errorf("update user error: %w", err)
	}

	return newTokens, nil
}

func (a *authService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	return a.jwTokens.ParseAccessToken(tokenString)
}
