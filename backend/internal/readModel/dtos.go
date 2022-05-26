package readModel

import (
	"time"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationFullDTO struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Avatar            string    `json:"avatar"`
	CreatedAt         time.Time `json:"created_at"`
	Type              string    `json:"type"`
	HasJoined         bool      `json:"joined,omitempty"`
	ParticipantsCount int64     `json:"participants_count,omitempty"`
	IsOwner           bool      `json:"is_owner,omitempty"`
}

type UserDTO struct {
	ID     uuid.UUID `json:"id"`
	Avatar string    `json:"avatar"`
	Name   string    `json:"name"`
}

type ContactDTO struct {
	ID     uuid.UUID `json:"id"`
	Avatar string    `json:"avatar"`
	Name   string    `json:"name"`
}

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text,omitempty"`
	Type      string    `json:"type"`
	User      *UserDTO  `json:"user"`
	IsInbound bool      `json:"is_inbound,omitempty"`
	NewName   string    `json:"new_name,omitempty"`
}
