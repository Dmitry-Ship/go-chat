package domain

import (
	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid" json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	IsPrivate bool      `json:"is_private"`
}

func NewConversation(id uuid.UUID, name string, isPrivate bool) *Conversation {
	return &Conversation{
		ID:        id,
		Name:      name,
		Avatar:    string(name[0]),
		IsPrivate: isPrivate,
	}
}
