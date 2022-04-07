package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID  `gorm:"type:uuid" json:"id"`
	ConversationID uuid.UUID  `gorm:"type:uuid" json:"conversation_id"`
	UserID         *uuid.UUID `gorm:"type:uuid" json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	Text           string     `json:"text"`
	Type           int        `json:"type"`
}

func NewUserMessage(text string, conversationId uuid.UUID, userID uuid.UUID) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Text:           text,
		Type:           0,
		UserID:         &userID,
	}
}

func NewSystemMessage(text string, conversationId uuid.UUID) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Text:           text,
		Type:           1,
	}
}
