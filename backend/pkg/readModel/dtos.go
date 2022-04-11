package readModel

import (
	"GitHub/go-chat/backend/domain"
	"time"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationDTOFull struct {
	Conversation ConversationDTO `json:"conversation"`
	HasJoined    bool            `json:"joined"`
}

type UserDTO struct {
	ID     uuid.UUID `json:"id"`
	Avatar string    `json:"avatar"`
	Name   string    `json:"name"`
}

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	User      *UserDTO  `json:"user,omitempty"`
	IsInbound bool      `json:"is_inbound,omitempty"`
}

func ToMessageDTOFromDomain(message *domain.Message) *MessageDTO {
	return &MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Text:      message.Text,
		Type:      message.Type,
	}
}
