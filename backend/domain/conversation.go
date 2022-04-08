package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID
	Name      string
	Avatar    string
	IsPrivate bool
	CreatedAt time.Time
}

func NewConversation(id uuid.UUID, name string, isPrivate bool) *Conversation {
	return &Conversation{
		ID:        id,
		Name:      name,
		Avatar:    string(name[0]),
		IsPrivate: isPrivate,
		CreatedAt: time.Now(),
	}
}
