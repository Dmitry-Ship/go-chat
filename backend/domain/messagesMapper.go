package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	User      *UserDTO  `json:"user,omitempty"`
	IsInbound bool      `json:"is_inbound,omitempty"`
}

type MessageDAO struct {
	ID             uuid.UUID  `gorm:"type:uuid"`
	ConversationID uuid.UUID  `gorm:"type:uuid"`
	UserID         *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
	Text           string
	Type           uint8
}

func (MessageDAO) TableName() string {
	return "messages"
}

func ToMessageDTO(message *MessageDAO, user *UserDAO, requestUserID uuid.UUID) *MessageDTO {
	messageDTO := MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Text:      message.Text,
		Type:      "system",
	}

	if message.Type == 0 {
		messageDTO.User = ToUserDTO(user)
		messageDTO.IsInbound = user.ID != requestUserID
		messageDTO.Type = "user"
	}

	return &messageDTO
}

func ToMessageDTOFromDomain(message *Message) *MessageDTO {
	return &MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Text:      message.Text,
		Type:      message.Type,
	}
}

func ToMessageDAO(message *Message) *MessageDAO {
	var messageType uint8

	switch message.Type {
	case "user":
		messageType = 0
	case "system":
		messageType = 1
	}

	return &MessageDAO{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		CreatedAt:      message.CreatedAt,
		Text:           message.Text,
		Type:           messageType,
	}
}
