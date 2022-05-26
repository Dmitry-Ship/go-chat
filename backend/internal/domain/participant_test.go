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

func TestLeaveGroupConversation(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	participant := NewParticipant(conversationID, userID)

	err := participant.LeaveGroupConversation(conversationID)

	assert.Nil(t, err)
	assert.Equal(t, participant.IsActive, false)
	assert.Equal(t, participant.events[len(participant.events)-1], NewGroupConversationLeft(conversationID, userID))
}
