package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("test", "12345678")

	assert.NotNil(t, user.ID)
	assert.NotNil(t, user.Avatar)
	assert.Equal(t, user.Password, "12345678")
	assert.Equal(t, user.Name, "test")
	assert.Nil(t, err)
}

func TestNewUserEmptyUsername(t *testing.T) {
	_, err := NewUser("", "123")

	assert.Equal(t, err.Error(), "username is empty")
}

func TestNewUserEmptyPassword(t *testing.T) {
	_, err := NewUser("test", "")

	assert.Equal(t, err.Error(), "password is empty")
}

func TestNewUserShortPassword(t *testing.T) {
	_, err := NewUser("test", "123")

	assert.Equal(t, err.Error(), "password is too short")
}

func TestUser_SetRefreshToken(t *testing.T) {
	user, err := NewUser("test", "12345678")

	user.SetRefreshToken("test")

	assert.Nil(t, err)
	assert.Equal(t, user.RefreshToken, "test")
}
