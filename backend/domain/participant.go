package domain

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID             uuid.UUID `gorm:"type:uuid" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid" json:"conversation_id"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	return &Participant{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userId,
		CreatedAt:      time.Now(),
	}
}
