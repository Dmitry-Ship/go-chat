package domain

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantAggregate struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	CreatedAt      time.Time
}

func NewParticipant(conversationId uuid.UUID, userId uuid.UUID) *ParticipantAggregate {
	return &ParticipantAggregate{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userId,
		CreatedAt:      time.Now(),
	}
}
