package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser("test", "123")
	assert.NotNil(t, user.ID)
	assert.NotNil(t, user.Avatar)
	assert.Equal(t, user.Password, "123")
	assert.Equal(t, user.Name, "test")
}

func TestUser_SetRefreshToken(t *testing.T) {
	user := NewUser("test", "123")
	user.SetRefreshToken("test")
	assert.Equal(t, user.RefreshToken, "test")
}
