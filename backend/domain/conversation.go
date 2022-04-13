package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID
	Name      string
	Avatar    string
	Type      string
	CreatedAt time.Time
}

func NewPublicConversation(id uuid.UUID, name string) *Conversation {
	return &Conversation{
		ID:        id,
		Name:      name,
		Avatar:    string(name[0]),
		Type:      "public",
		CreatedAt: time.Now(),
	}
}
