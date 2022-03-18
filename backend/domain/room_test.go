package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRoom(t *testing.T) {
	name := "test"
	roomId := uuid.New()
	room := NewRoom(roomId, name)
	assert.Equal(t, room.ID, roomId)
	assert.Equal(t, name, room.Name)
}
