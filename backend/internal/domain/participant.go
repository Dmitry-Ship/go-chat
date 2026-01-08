package domain

import (
	"context"

	"github.com/google/uuid"
)

type ParticipantRepository interface {
	Store(ctx context.Context, participant *Participant) error
	Update(ctx context.Context, participant *Participant) error
	GetByConversationIDAndUserID(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*Participant, error)
	GetIDsByConversationID(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
	GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type Participant struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	IsActive       bool
}

func NewParticipant(participantID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID) *Participant {
	return &Participant{
		ID:             participantID,
		ConversationID: conversationID,
		UserID:         userID,
		IsActive:       true,
	}
}
