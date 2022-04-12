package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTextMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	message := NewTextMessage(conversationID, userID, "content")

	assert.Equal(t, "content", message.Data.Text)
	assert.Equal(t, 0, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
}
