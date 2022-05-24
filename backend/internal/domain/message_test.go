package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTextMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message, err := NewTextMessage(conversationID, userID, "content")

	assert.Equal(t, "content", message.Data.Text)
	assert.NotNil(t, message.Data.ID)
	assert.Equal(t, MessageTypeText, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.events[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
	assert.Nil(t, err)
}

func TestNewTextMessageEmptyText(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	_, err := NewTextMessage(conversationID, userID, "")

	assert.Equal(t, err.Error(), "text is empty")
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
	assert.Equal(t, message.events[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewLeftConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := NewLeftConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeLeftConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.events[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewJoinedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := NewJoinedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeJoinedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.events[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewInvitedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := NewInvitedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeInvitedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.events[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}
