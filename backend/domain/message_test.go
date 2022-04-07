package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	message := NewMessage("content", "message", conversationID, userID)

	assert.Equal(t, "content", message.Content)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
}
