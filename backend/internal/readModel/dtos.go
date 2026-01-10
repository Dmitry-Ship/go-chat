package readModel

import (
	"time"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Avatar      string     `json:"avatar"`
	Type        string     `json:"type"`
	LastMessage MessageDTO `json:"last_message"`
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
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Text           string    `json:"text,omitempty"`
	Type           string    `json:"type"`
	User           UserDTO   `json:"user"`
	IsInbound      bool      `json:"is_inbound,omitempty"`
	ConversationId uuid.UUID `json:"conversation_id"`
}

type RawMessageDTO struct {
	ID             uuid.UUID
	Type           uint8
	CreatedAt      time.Time
	ConversationID uuid.UUID
	Content        string
	UserID         uuid.UUID
	UserName       string
	UserAvatar     string
}

type RawLastMessageDTO struct {
	MessageID         uuid.UUID
	MessageCreatedAt  time.Time
	MessageContent    string
	MessageType       int32
	MessageUserID     uuid.UUID
	MessageUserName   string
	MessageUserAvatar string
	ConversationID    uuid.UUID
}
