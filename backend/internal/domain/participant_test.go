package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewParticipant(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	participant := NewParticipant(conversationID, userID)

	assert.Equal(t, conversationID, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
}

func TestNewJoinedParticipant(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	participant := NewJoinedParticipant(conversationID, userID)

	assert.Equal(t, conversationID, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.events[len(participant.events)-1], NewPublicConversationJoined(conversationID, userID))
}

func TestNewInvitedParticipant(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()

	participant := NewInvitedParticipant(conversationID, userID)

	assert.Equal(t, conversationID, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.NotNil(t, participant.ID)
	assert.NotNil(t, participant.CreatedAt)
	assert.Equal(t, participant.IsActive, true)
	assert.Equal(t, participant.events[len(participant.events)-1], NewPublicConversationInvited(conversationID, userID))
}

func TestLeavePublicConversation(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	participant := NewParticipant(conversationID, userID)

	err := participant.LeavePublicConversation(conversationID)

	assert.Nil(t, err)
	assert.Equal(t, participant.events[len(participant.events)-1], NewPublicConversationLeft(conversationID, userID))
}
