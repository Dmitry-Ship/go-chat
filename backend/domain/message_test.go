package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	message := NewUserMessage("content", conversationID, userID)

	assert.Equal(t, "content", message.Text)
	assert.Equal(t, 0, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, *message.UserID)
	assert.NotNil(t, message.ID)
}

func TestNewSystemMessage(t *testing.T) {
	conversationID := uuid.New()
	message := NewSystemMessage("content", conversationID)

	assert.Equal(t, "content", message.Text)
	assert.Equal(t, 1, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.NotNil(t, message.ID)
}
