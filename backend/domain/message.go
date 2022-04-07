package domain

import (
	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid" json:"conversation_id"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"-"`
	Content        string    `json:"content"`
	Type           string    `json:"type"`
}

func NewMessage(content string, messageType string, conversationId uuid.UUID, userID uuid.UUID) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		Content:        content,
		Type:           messageType,
		UserID:         userID,
	}
}
