package domain

import (
	"github.com/google/uuid"
)

type Room struct {
	Name      string    `json:"name"`
	ID        uuid.UUID `json:"id" gorm:"type:uuid"`
	IsPrivate bool      `json:"is_private"`
}

func NewRoom(id uuid.UUID, name string, isPrivate bool) *Room {
	return &Room{
		ID:        id,
		Name:      name,
		IsPrivate: isPrivate,
	}
}
