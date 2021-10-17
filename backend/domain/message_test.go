package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewChatMessage(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	message := NewChatMessage("content", "message", roomID, userID)

	assert.Equal(t, "content", message.Content)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, roomID, message.RoomId)
	assert.Equal(t, userID, message.UserId)
	assert.NotNil(t, message.CreatedAt)
	assert.NotNil(t, message.Id)
}
