package domain

import (
	"github.com/google/uuid"
)

type Participant struct {
	ID             uuid.UUID `gorm:"type:uuid" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid" json:"conversation_id"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
}

func NewParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	return &Participant{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userId,
	}
}
