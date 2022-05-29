package domain

import (
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	GenericRepository[*Participant]
	GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*Participant, error)
}

type Participant struct {
	aggregate
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	IsActive       bool
}

type ParticipantLeaver interface {
	LeaveGroupConversation(conversationID uuid.UUID) error
}

func NewParticipant(participantID uuid.UUID, conversationID uuid.UUID, userID uuid.UUID) *Participant {
	return &Participant{
		ID:             participantID,
		ConversationID: conversationID,
		UserID:         userID,
		IsActive:       true,
	}
}
