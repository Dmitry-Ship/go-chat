package domain

import (
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	Store(participant *Participant) error
	Update(participant *Participant) error
	GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*Participant, error)
	GetIDsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error)
	GetConversationIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error)
}

type Participant struct {
	aggregate
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
