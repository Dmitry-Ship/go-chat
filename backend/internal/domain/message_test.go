package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTextMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message, err := newTextMessage(conversationID, userID, "content")

	assert.Equal(t, "content", message.Data.Text)
	assert.NotNil(t, message.Data.ID)
	assert.Equal(t, MessageTypeText, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
	assert.Nil(t, err)
}

func TestNewTextMessageEmptyText(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	_, err := newTextMessage(conversationID, userID, "")

	assert.Equal(t, err.Error(), "text is empty")
}

func TestNewTextMessageTooLong(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	maxTextLength := 1000
	text := ""

	for i := 0; i < maxTextLength+1; i++ {
		text += "a"
	}

	_, err := newTextMessage(conversationID, userID, text)

	assert.Equal(t, err.Error(), "text is too long")
}

func TestNewConversationRenamedMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newConversationRenamedMessage(conversationID, userID, "new name")

	assert.Equal(t, "new name", message.Data.NewName)
	assert.NotNil(t, message.Data.ID)
	assert.Equal(t, MessageTypeRenamedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewLeftConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newLeftConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeLeftConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewJoinedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newJoinedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeJoinedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewInvitedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newInvitedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeInvitedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.events)-1], NewMessageSent(conversationID, message.ID, userID))
}
