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

	assert.Equal(t, conversation.Conversation.ID, conversationID)
	assert.Equal(t, len(conversation.Participants), 2)
	assert.Equal(t, conversation.Type, ConversationTypeDirect)
	assert.Equal(t, true, conversation.IsActive)
	assert.Equal(t, conversation.GetEvents()[len(conversation.GetEvents())-1], newDirectConversationCreatedEvent(conversationID, []uuid.UUID{to, from}))
	assert.Nil(t, err)
}

func TestNewDirectConversationWithOneself(t *testing.T) {
	to := uuid.New()
	conversationID := uuid.New()

	_, err := NewDirectConversation(conversationID, to, to)

	assert.Equal(t, err.Error(), "cannot chat with yourself")
}

func TestSendDirectTextMessage(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()
	messageID := uuid.New()
	conversation, _ := NewDirectConversation(conversationID, to, from)
	text := "Hello world"

	message, err := conversation.SendTextMessage(messageID, text, &conversation.Participants[0])

	assert.Nil(t, err)
	assert.Equal(t, message.Content.String(), text)
	assert.Equal(t, message.UserID, to)
}

func TestSendDirectTextMessageNotAMember(t *testing.T) {
	to := uuid.New()
	from := uuid.New()
	conversationID := uuid.New()
	messageID := uuid.New()
	conversation, _ := NewDirectConversation(conversationID, to, from)
	text := "Hello world"
	participant := NewParticipant(uuid.New(), conversationID, uuid.New())

	_, err := conversation.SendTextMessage(messageID, text, participant)

	assert.Equal(t, err.Error(), "user is not participant")
}
