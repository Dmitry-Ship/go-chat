package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid" json:"conversation_id"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	Text           string    `json:"text"`
	Type           string    `json:"type"`
}

func NewMessage(text string, messageType string, conversationId uuid.UUID, userID uuid.UUID) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Text:           text,
		Type:           messageType,
		UserID:         userID,
	}
}
