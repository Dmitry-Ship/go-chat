package domain

import (
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	Store(participant *Participant) error
	GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*Participant, error)
	Update(participant *Participant) error
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

func NewParticipant(conversationID uuid.UUID, userID uuid.UUID) *Participant {
	return &Participant{
		ID:             uuid.New(),
		ConversationID: conversationID,
		UserID:         userID,
		IsActive:       true,
	}
}
