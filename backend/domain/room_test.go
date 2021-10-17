package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoom(t *testing.T) {
	name := "test"
	room := NewRoom(name)
	assert.NotNil(t, room.Id)
	assert.Equal(t, name, room.Name)
}
