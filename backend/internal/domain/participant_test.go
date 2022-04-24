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
}
