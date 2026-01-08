package services

import (
	"errors"
	"time"

	"GitHub/go-chat/backend/internal/config"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type tokenClaims struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type tokens struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	RefreshTokenExpiration time.Duration
	AccessTokenExpiration  time.Duration
}

type jwTokens struct {
	config config.Auth
}

type JWTokens interface {
	ParseAccessToken(accessTokenString string) (uuid.UUID, error)
	ParseRefreshToken(refreshTokenString string) (uuid.UUID, error)
	CreateTokens(userid uuid.UUID) (tokens, error)
}

func NewJWTokens(config config.Auth) *jwTokens {
	return &jwTokens{
		config: config,
	}
}

func (a *jwTokens) createAccessToken(userid uuid.UUID) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.config.AccessToken.TTL).Unix(),
		},
		UserID: userid,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(a.config.AccessToken.Secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *jwTokens) createRefreshToken(userid uuid.UUID) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.config.RefreshToken.TTL).Unix(),
		},
		UserID: userid,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(a.config.RefreshToken.Secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *jwTokens) CreateTokens(userid uuid.UUID) (tokens, error) {
	var newTokens tokens

	accessToken, err := a.createAccessToken(userid)

	if err != nil {
		return newTokens, err
	}

	refreshToken, err := a.createRefreshToken(userid)

	if err != nil {
		return newTokens, err
	}

	newTokens.AccessToken = accessToken
	newTokens.RefreshToken = refreshToken
	newTokens.AccessTokenExpiration = a.config.AccessToken.TTL
	newTokens.RefreshTokenExpiration = a.config.RefreshToken.TTL

	return newTokens, nil
}

func (a *jwTokens) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.config.AccessToken.Secret), nil
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

func (a *jwTokens) ParseRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.config.RefreshToken.Secret), nil
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
		userID := token.Claims.(*tokenClaims).UserID

		return userID, nil
	}

	return uuid.Nil, err
}
