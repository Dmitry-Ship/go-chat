package readModel

import (
	"time"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationFullDTO struct {
	Conversation ConversationDTO `json:"conversation"`
	HasJoined    bool            `json:"joined"`
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
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	Text                string    `json:"text,omitempty"`
	Type                string    `json:"type"`
	User                *UserDTO  `json:"user,omitempty"`
	IsInbound           bool      `json:"is_inbound,omitempty"`
	NewConversationName string    `json:"new_name,omitempty"`
}
