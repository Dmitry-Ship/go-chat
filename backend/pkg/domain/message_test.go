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
	assert.NotNil(t, message.Data.ID)
	assert.Equal(t, MessageTypeText, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
}

func TestNewConversationRenamedMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	message := NewConversationRenamedMessage(conversationID, userID, "new name")

	assert.Equal(t, "new name", message.Data.NewName)
	assert.NotNil(t, message.Data.ID)
	assert.Equal(t, MessageTypeRenamedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
}

func TestNewLeftConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	message := NewLeftConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeLeftConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
}

func TestNewJoinedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	message := NewJoinedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeJoinedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
}
