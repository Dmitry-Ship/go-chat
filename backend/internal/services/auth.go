package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"errors"

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

func NewAuthService(users domain.UserRepository, jwTokens JWTokens) *authService {
	return &authService{
		users:    users,
		jwTokens: jwTokens,
	}
}

func (a *authService) Login(username string, password string) (tokens, error) {
	user, err := a.users.FindByUsername(username)

	if err != nil {
		return tokens{}, err
	}

	userPassword, err := domain.NewUserPassword(password, func(p []byte) ([]byte, error) {
		return []byte(p), nil
	})

	if err != nil {
		return tokens{}, err
	}

	if err = user.Password.Compare(userPassword, func(p1 []byte, p2 []byte) error {
		return bcrypt.CompareHashAndPassword(p1, p2)
	}); err != nil {
		return tokens{}, err
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, err
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Update(user); err != nil {
		return tokens{}, err
	}

	return newTokens, err
}

func (a *authService) Logout(userID uuid.UUID) error {
	user, err := a.users.GetByID(userID)

	if err != nil {
		return err
	}

	user.SetRefreshToken("")

	err = a.users.Update(user)

	return err
}

func (a *authService) SignUp(username string, password string) (tokens, error) {
	name, err := domain.NewUserName(username)

	if err != nil {
		return tokens{}, err
	}

	hashedPassword, err := domain.NewUserPassword(password, func(p []byte) ([]byte, error) {
		return bcrypt.GenerateFromPassword(p, 14)
	})

	if err != nil {
		return tokens{}, err
	}

	user := domain.NewUser(name, hashedPassword)

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, err
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Store(user); err != nil {
		return tokens{}, err
	}

	return newTokens, err
}

func (a *authService) RotateTokens(refreshTokenString string) (tokens, error) {
	userID, err := a.jwTokens.ParseRefreshToken(refreshTokenString)

	if err != nil {
		return tokens{}, err
	}

	user, err := a.users.GetByID(userID)

	if err != nil {
		return tokens{}, err
	}

	if user.RefreshToken != refreshTokenString {
		return tokens{}, errors.New("invalid token")
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, err
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	if err = a.users.Update(user); err != nil {
		return tokens{}, err
	}

	return newTokens, err
}

func (a *authService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	return a.jwTokens.ParseAccessToken(tokenString)
}
