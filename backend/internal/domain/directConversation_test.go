package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDirectConversation(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()

	conversation, err := NewDirectConversation(conversationID, to, from)

	assert.Equal(t, conversation.ID, conversationID)
	assert.Equal(t, to, conversation.Data.ToUser.UserID)
	assert.Equal(t, from, conversation.Data.FromUser.UserID)
	assert.Equal(t, conversationID, conversation.Data.FromUser.ConversationID)
	assert.Equal(t, conversationID, conversation.Data.ToUser.ConversationID)
	assert.NotNil(t, conversation.Data.FromUser.ID)
	assert.NotNil(t, conversation.Data.ToUser.ID)
	assert.Equal(t, conversation.Type, "direct")
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newDirectConversationCreatedEvent(conversationID, to, from))
	assert.Nil(t, err)
}

func TestNewDirectConversationWithOneself(t *testing.T) {
	to := uuid.New()
	conversationID := uuid.New()

	_, err := NewDirectConversation(conversationID, to, to)

	assert.Equal(t, err.Error(), "cannot chat with yourself")
}

func TestGetFromUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()

	conversation, _ := NewDirectConversation(conversationID, to, from)

	assert.Equal(t, from, conversation.GetFromUser().UserID)
}

func TestGetToUser(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()

	conversation, _ := NewDirectConversation(conversationID, to, from)

	assert.Equal(t, to, conversation.GetToUser().UserID)
}

func TestSendDirectTextMessage(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()

	conversation, _ := NewDirectConversation(conversationID, to, from)

	text := "Hello world"

	message, err := conversation.SendTextMessage(text, to)

	assert.Nil(t, err)
	assert.Equal(t, message.Data.Text, text)
	assert.Equal(t, message.UserID, to)
}

func TestSendDirectTextMessageNotAMember(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()

	conversation, _ := NewDirectConversation(conversationID, to, from)

	text := "Hello world"

	_, err := conversation.SendTextMessage(text, uuid.New())

	assert.Equal(t, err.Error(), "user is not participant")
}
