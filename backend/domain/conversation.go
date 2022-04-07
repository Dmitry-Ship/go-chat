package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid" json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
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
