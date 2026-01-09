package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewParticipant(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	participantID := uuid.New()

	participant := NewParticipant(participantID, conversationID, userID)

	assert.Equal(t, conversationID, participant.ConversationID)
	assert.Equal(t, userID, participant.UserID)
	assert.Equal(t, participantID, participant.ID)
}
