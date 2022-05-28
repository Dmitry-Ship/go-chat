package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTextMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	content, _ := newTextMessageContent("content")

	message := newTextMessage(conversationID, userID, content)

	assert.Equal(t, content, message.Content)
	assert.NotNil(t, message.ID)
	assert.Equal(t, MessageTypeText, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.GetEvents())-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewTextMessageContentEmpty(t *testing.T) {
	_, err := newTextMessageContent("")

	assert.Equal(t, err.Error(), "text is empty")
}

func TestNewTextMessageContentTooLong(t *testing.T) {
	maxTextLength := 1000
	text := ""

	for i := 0; i < maxTextLength+1; i++ {
		text += "a"
	}

	_, err := newTextMessageContent(text)

	assert.Equal(t, err.Error(), "text is too long")
}

func TestNewConversationRenamedMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	name := newRenamedMessageContent("new name")

	message := newConversationRenamedMessage(conversationID, userID, name)

	assert.Equal(t, name, message.Content)
	assert.NotNil(t, message.ID)
	assert.Equal(t, MessageTypeRenamedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.GetEvents())-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewLeftConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newLeftConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeLeftConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.GetEvents())-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewJoinedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newJoinedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeJoinedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.GetEvents())-1], NewMessageSent(conversationID, message.ID, userID))
}

func TestNewInvitedConversationMessage(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	message := newInvitedConversationMessage(conversationID, userID)

	assert.Equal(t, MessageTypeInvitedConversation, message.Type)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.UserID)
	assert.NotNil(t, message.ID)
	assert.Equal(t, message.GetEvents()[len(message.GetEvents())-1], NewMessageSent(conversationID, message.ID, userID))
}
