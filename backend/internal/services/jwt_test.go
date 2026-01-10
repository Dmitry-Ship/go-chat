package services

import (
	"testing"
	"time"

	"GitHub/go-chat/backend/internal/config"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func getTestConfig() config.Auth {
	return config.Auth{
		AccessToken: config.Token{
			Secret: "test-access-secret-key-32-bytes-long",
			TTL:    15 * time.Minute,
		},
		RefreshToken: config.Token{
			Secret: "test-refresh-secret-key-32-bytes",
			TTL:    7 * 24 * time.Hour,
		},
	}
}

func TestNewJWTokens(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())

	assert.NotNil(t, jwtService)
	assert.IsType(t, &jwTokens{}, jwtService)
}

func TestJWTokens_CreateTokens(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())
	userID := uuid.New()

	tokens, err := jwtService.CreateTokens(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, 15*time.Minute, tokens.AccessTokenExpiration)
	assert.Equal(t, 7*24*time.Hour, tokens.RefreshTokenExpiration)
}

func TestJWTokens_ParseAccessToken(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())
	userID := uuid.New()

	tokens, err := jwtService.CreateTokens(userID)
	assert.NoError(t, err)

	parsedUserID, err := jwtService.ParseAccessToken(tokens.AccessToken)

	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestJWTokens_ParseAccessToken_Invalid(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())
	userID := uuid.New()

	tokens, err := jwtService.CreateTokens(userID)
	assert.NoError(t, err)

	_, err = jwtService.ParseAccessToken(tokens.AccessToken + "modified")

	assert.Error(t, err)
}

func TestJWTokens_ParseAccessToken_Malformed(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())

	_, err := jwtService.ParseAccessToken("a.b.c")

	assert.Error(t, err)
}

func TestJWTokens_ParseAccessToken_WrongSecret(t *testing.T) {
	jwtService1 := NewJWTokens(getTestConfig())
	jwtService2 := NewJWTokens(config.Auth{
		AccessToken: config.Token{
			Secret: "different-secret-key-32-bytes-long",
			TTL:    15 * time.Minute,
		},
		RefreshToken: config.Token{
			Secret: "test-refresh-secret-key-32-bytes",
			TTL:    7 * 24 * time.Hour,
		},
	})

	userID := uuid.New()
	tokens, _ := jwtService1.CreateTokens(userID)

	_, err := jwtService2.ParseAccessToken(tokens.AccessToken)

	assert.Error(t, err)
}

func TestJWTokens_ParseRefreshToken(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())
	userID := uuid.New()

	tokens, err := jwtService.CreateTokens(userID)
	assert.NoError(t, err)

	parsedUserID, err := jwtService.ParseRefreshToken(tokens.RefreshToken)

	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestJWTokens_ParseRefreshToken_Invalid(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())

	_, err := jwtService.ParseRefreshToken("invalid-token")

	assert.Error(t, err)
}

func TestJWTokens_ParseRefreshToken_Malformed(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())

	_, err := jwtService.ParseRefreshToken("not.a.valid.jwt.token")

	assert.Error(t, err)
}

func TestJWTokens_ParseRefreshToken_WrongSecret(t *testing.T) {
	jwtService1 := NewJWTokens(getTestConfig())
	jwtService2 := NewJWTokens(config.Auth{
		AccessToken: config.Token{
			Secret: "test-access-secret-key-32-bytes-long",
			TTL:    15 * time.Minute,
		},
		RefreshToken: config.Token{
			Secret: "different-secret-key-32-bytes",
			TTL:    7 * 24 * time.Hour,
		},
	})

	userID := uuid.New()
	tokens, _ := jwtService1.CreateTokens(userID)

	_, err := jwtService2.ParseRefreshToken(tokens.RefreshToken)

	assert.Error(t, err)
}

func TestJWTokens_CreateTokens_DifferentUserIDs(t *testing.T) {
	jwtService := NewJWTokens(getTestConfig())
	userID1 := uuid.New()
	userID2 := uuid.New()

	tokens1, err1 := jwtService.CreateTokens(userID1)
	tokens2, err2 := jwtService.CreateTokens(userID2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, tokens1.AccessToken, tokens2.AccessToken)
	assert.NotEqual(t, tokens1.RefreshToken, tokens2.RefreshToken)

	parsedID1, _ := jwtService.ParseAccessToken(tokens1.AccessToken)
	parsedID2, _ := jwtService.ParseAccessToken(tokens2.AccessToken)

	assert.Equal(t, userID1, parsedID1)
	assert.Equal(t, userID2, parsedID2)
}

func TestJWTokens_ParseAccessToken_Expired(t *testing.T) {
	jwtService := NewJWTokens(config.Auth{
		AccessToken: config.Token{
			Secret: "test-access-secret-key-32-bytes-long",
			TTL:    -1 * time.Second,
		},
		RefreshToken: config.Token{
			Secret: "test-refresh-secret-key-32-bytes",
			TTL:    7 * 24 * time.Hour,
		},
	})
	userID := uuid.New()

	tokens, err := jwtService.CreateTokens(userID)
	assert.NoError(t, err)

	_, err = jwtService.ParseAccessToken(tokens.AccessToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}
