package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type tokenClaims struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authService struct {
	users                  domain.UserCommandRepository
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

type AuthService interface {
	Login(username string, password string) (Tokens, error)
	Logout(userId uuid.UUID) error
	SignUp(username string, password string) (Tokens, error)
	RotateTokens(refreshTokenString string) (Tokens, error)
	GetAccessTokenExpiration() time.Duration
	GetRefreshTokenExpiration() time.Duration
	ParseAccessToken(accessTokenString string) (uuid.UUID, error)
}

func NewAuthService(users domain.UserCommandRepository) *authService {
	return &authService{
		users:                  users,
		refreshTokenExpiration: 24 * 90 * time.Hour,
		accessTokenExpiration:  10 * time.Minute,
	}
}

func (a *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *authService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *authService) Login(username string, password string) (Tokens, error) {
	user, err := a.users.FindByUsername(username)

	if err != nil {
		return Tokens{}, err
	}

	if !a.checkPasswordHash(password, user.Password) {
		return Tokens{}, errors.New("password is incorrect")
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

func (a *authService) SignUp(username string, password string) (Tokens, error) {
	hashedPassword, err := a.hashPassword(password)

	if err != nil {
		return Tokens{}, err
	}

	user := domain.NewUser(username, hashedPassword)

	tokens, err := a.createTokens(user.ID)

	if err != nil {
		return Tokens{}, err
	}

	user.SetRefreshToken(tokens.RefreshToken)

	err = a.users.Store(user)

	if err != nil {
		return Tokens{}, err
	}

	return tokens, err
}

func (a *authService) createAccessToken(userid uuid.UUID) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.accessTokenExpiration).Unix(),
		},
		UserID: userid,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) createRefreshToken(userid uuid.UUID) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.refreshTokenExpiration).Unix(),
		},
		UserID: userid,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) createTokens(userid uuid.UUID) (Tokens, error) {
	var tokens Tokens

	accessToken, err := a.createAccessToken(userid)

	if err != nil {
		return tokens, err
	}

	refreshToken, err := a.createRefreshToken(userid)

	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (a *authService) createAndSetTokens(user *domain.User) (Tokens, error) {
	tokens, err := a.createTokens(user.ID)

	if err != nil {
		return Tokens{}, err
	}

	user.SetRefreshToken(tokens.RefreshToken)

	err = a.users.Update(user)

	if err != nil {
		return Tokens{}, err
	}

	return tokens, err
}

func (a *authService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	if token.Valid {
		return token.Claims.(*tokenClaims).UserID, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return uuid.Nil, errors.New("invalid token")
		}

		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return uuid.Nil, errors.New("token expired")
		}

		return uuid.Nil, err
	}

	return uuid.Nil, err
}

func (a *authService) parseRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return uuid.Nil, errors.New("invalid token")
			}

			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return uuid.Nil, errors.New("token expired")
			}

			return uuid.Nil, err
		}

		return uuid.Nil, err
	}

	if token.Valid {
		userId := token.Claims.(*tokenClaims).UserID

		user, err := a.users.GetByID(userId)

		if err != nil {
			return uuid.Nil, err
		}

		if user.RefreshToken != tokenString {
			return uuid.Nil, errors.New("invalid token")
		}

		return userId, nil
	}

	return uuid.Nil, err
}

func (a *authService) RotateTokens(refreshTokenString string) (Tokens, error) {
	userId, err := a.parseRefreshToken(refreshTokenString)

	if err != nil {
		return Tokens{}, err
	}

	user, err := a.users.GetByID(userId)

	if err != nil {
		return Tokens{}, err
	}

	return a.createAndSetTokens(user)
}

func (a *authService) GetAccessTokenExpiration() time.Duration {
	return a.accessTokenExpiration
}

func (a *authService) GetRefreshTokenExpiration() time.Duration {
	return a.refreshTokenExpiration
}
