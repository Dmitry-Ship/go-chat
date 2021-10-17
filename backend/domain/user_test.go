package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser()
	assert.NotNil(t, user.Id)
	assert.NotNil(t, user.Avatar)
	assert.NotNil(t, user.Name)
}
