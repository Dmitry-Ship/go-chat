package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser("test", "123")
	assert.NotNil(t, user.Id)
	assert.NotNil(t, user.Avatar)
	assert.Equal(t, user.Password, "123")
	assert.Equal(t, user.Name, "test")
}
