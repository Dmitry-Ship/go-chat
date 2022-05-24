package services

import (
	"errors"
	"os"
	"time"

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
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

type JWTokens interface {
	ParseAccessToken(accessTokenString string) (uuid.UUID, error)
	ParseRefreshToken(refreshTokenString string) (uuid.UUID, error)
	CreateTokens(userid uuid.UUID) (tokens, error)
}

func NewJWTokens() *jwTokens {
	return &jwTokens{
		refreshTokenExpiration: 24 * 90 * time.Hour,
		accessTokenExpiration:  10 * time.Minute,
	}
}

func (a *jwTokens) createAccessToken(userid uuid.UUID) (string, error) {
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

func (a *jwTokens) createRefreshToken(userid uuid.UUID) (string, error) {
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
	newTokens.AccessTokenExpiration = a.accessTokenExpiration
	newTokens.RefreshTokenExpiration = a.refreshTokenExpiration

	return newTokens, nil
}

func (a *jwTokens) ParseAccessToken(tokenString string) (uuid.UUID, error) {
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

func (a *jwTokens) ParseRefreshToken(tokenString string) (uuid.UUID, error) {
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
		userID := token.Claims.(*tokenClaims).UserID

		return userID, nil
	}

	return uuid.Nil, err
}
