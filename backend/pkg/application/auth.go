package application

import (
	"GitHub/go-chat/backend/domain"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const RefreshTokenExpiration = 24 * 7 * time.Hour
const AccessTokenExpiration = 30 * time.Second

type tokenClaims struct {
	UserId uuid.UUID
	jwt.StandardClaims
}

type authService struct {
	users domain.UserRepository
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Login(username string, password string) (Tokens, error)
	Logout(userId uuid.UUID) error
	SignUp(username string, password string) (Tokens, error)
	RefreshAccessToken(refreshTokenString string) (string, error)
	GetUser(userId uuid.UUID) (*domain.User, error)
}

func NewAuthService(users domain.UserRepository) *authService {
	return &authService{
		users: users,
	}
}

func (s *authService) GetUser(userId uuid.UUID) (*domain.User, error) {
	return s.users.FindByID(userId)
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

	newTokens, err := a.createTokens(user.ID)

	if err != nil {
		return newTokens, err
	}

	a.users.StoreRefreshToken(user.ID, newTokens.RefreshToken)

	return newTokens, nil
}

func (a *authService) Logout(userId uuid.UUID) error {
	err := a.users.DeleteRefreshToken(userId)

	if err != nil {
		return err
	}

	return nil
}

func (a *authService) SignUp(username string, password string) (Tokens, error) {
	hashedPassword, err := a.hashPassword(password)

	if err != nil {
		return Tokens{}, err
	}

	user := domain.NewUser(username, hashedPassword)

	err = a.users.Store(user)

	if err != nil {
		return Tokens{}, err
	}

	newTokens, err := a.createTokens(user.ID)

	if err != nil {
		return newTokens, err
	}

	err = a.users.StoreRefreshToken(user.ID, newTokens.RefreshToken)

	if err != nil {
		return newTokens, err
	}

	return newTokens, nil
}

func (a *authService) createAccessToken(userid uuid.UUID) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpiration).Unix(),
		},
		UserId: userid,
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
			ExpiresAt: time.Now().Add(RefreshTokenExpiration).Unix(),
		},
		UserId: userid,
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

func (a *authService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	if token.Valid {
		return token.Claims.(*tokenClaims).UserId, nil
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
		userId := token.Claims.(*tokenClaims).UserId
		refreshToken, err := a.users.GetRefreshTokenByUserId(userId)

		if err != nil {
			return uuid.Nil, err
		}

		if refreshToken != tokenString {
			return uuid.Nil, errors.New("invalid token")
		}

		return userId, nil
	}

	return uuid.Nil, err
}

func (a *authService) RefreshAccessToken(refreshTokenString string) (string, error) {
	userId, err := a.parseRefreshToken(refreshTokenString)

	if err != nil {
		return "", err
	}

	return a.createAccessToken(userId)
}
