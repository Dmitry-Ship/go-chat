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

func TestNewUserLongUsername(t *testing.T) {
	name := ""

	for i := 0; i < 101; i++ {
		name += "a"
	}

	_, err := NewUser(name, "123")

	assert.Equal(t, err.Error(), "username is too long")
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
