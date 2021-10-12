package domain

import (
	"github.com/google/uuid"
)

type Room struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

func NewRoom(name string) *Room {
	return &Room{
		Id:   uuid.New(),
		Name: name,
	}
}
