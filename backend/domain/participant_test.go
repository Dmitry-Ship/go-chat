package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewParticipant(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	participant := NewParticipant(roomID, userID)
	assert.Equal(t, roomID, participant.RoomId)
	assert.Equal(t, userID, participant.UserId)
	assert.NotNil(t, participant.CreatedAt)
	assert.NotNil(t, participant.ID)
}
