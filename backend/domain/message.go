package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         *uuid.UUID
	CreatedAt      time.Time
	Text           string
	Type           string
}

func NewUserMessage(text string, conversationId uuid.UUID, userID uuid.UUID) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Text:           text,
		Type:           "user",
		UserID:         &userID,
	}
}

func NewSystemMessage(text string, conversationId uuid.UUID) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationId,
		CreatedAt:      time.Now(),
		Text:           text,
		Type:           "system",
	}
}
