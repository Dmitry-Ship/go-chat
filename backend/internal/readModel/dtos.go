package readModel

import (
	"time"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Avatar     string    `json:"avatar"`
	CreatedAt  time.Time `json:"created_at"`
	Type       uint8     `json:"-"`
	UserID     uuid.UUID `json:"-"`
	UserAvatar string    `json:"-"`
	UserName   string    `json:"-"`
}

type ConversationFullDTO struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Avatar            string    `json:"avatar"`
	CreatedAt         time.Time `json:"created_at"`
	UserID            uuid.UUID `json:"-"`
	UserAvatar        string    `json:"-"`
	UserName          string    `json:"-"`
	PersistenceType   uint8     `json:"-"`
	Type              string    `json:"type"`
	HasJoined         bool      `json:"joined"`
	ParticipantsCount int64     `json:"participants_count"`
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
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Text            string    `json:"text,omitempty"`
	PersistenceType uint8     `json:"-"`
	Type            string    `json:"type"`
	UserID          uuid.UUID `json:"-"`
	UserName        string    `json:"-"`
	UserAvatar      string    `json:"-"`
	User            *UserDTO  `json:"user"`
	IsInbound       bool      `json:"is_inbound,omitempty"`
	NewName         string    `json:"new_name,omitempty"`
}
