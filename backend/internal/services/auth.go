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
	Logout(userId uuid.UUID) error
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return tokens{}, errors.New("password is incorrect")
	}

	return a.createAndSetTokens(user)
}

func (a *authService) Logout(userId uuid.UUID) error {
	user, err := a.users.GetByID(userId)

	if err != nil {
		return err
	}

	user.SetRefreshToken("")

	err = a.users.Update(user)

	return err
}

func (a *authService) SignUp(username string, password string) (tokens, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return tokens{}, err
	}

	hashedPassword := string(bytes)

	user, err := domain.NewUser(username, hashedPassword)

	if err != nil {
		return tokens{}, err
	}

	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, err
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	err = a.users.Store(user)

	if err != nil {
		return tokens{}, err
	}

	return newTokens, err
}

func (a *authService) RotateTokens(refreshTokenString string) (tokens, error) {
	userId, err := a.jwTokens.ParseRefreshToken(refreshTokenString)

	if err != nil {
		return tokens{}, err
	}

	user, err := a.users.GetByID(userId)

	if err != nil {
		return tokens{}, err
	}

	if user.RefreshToken != refreshTokenString {
		return tokens{}, errors.New("invalid token")
	}

	return a.createAndSetTokens(user)
}

func (a *authService) createAndSetTokens(user *domain.User) (tokens, error) {
	newTokens, err := a.jwTokens.CreateTokens(user.ID)

	if err != nil {
		return tokens{}, err
	}

	user.SetRefreshToken(newTokens.RefreshToken)

	err = a.users.Update(user)

	if err != nil {
		return tokens{}, err
	}

	return newTokens, err
}

func (a *authService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	return a.jwTokens.ParseAccessToken(tokenString)
}
