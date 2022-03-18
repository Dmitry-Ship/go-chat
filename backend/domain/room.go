package domain

import (
	"github.com/google/uuid"
)

type Room struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id" gorm:"type:uuid"`
}

func NewRoom(id uuid.UUID, name string) *Room {
	return &Room{
		ID:   id,
		Name: name,
	}
}
