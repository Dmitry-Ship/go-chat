package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageDTO struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UserId    *uuid.UUID `json:"user_id"`
	Text      string     `json:"text"`
	Type      string     `json:"type"`
	User      *UserDTO   `json:"user,omitempty"`
	IsInbound bool       `json:"is_inbound,omitempty"`
}

type MessagePersistence struct {
	ID             uuid.UUID  `gorm:"type:uuid"`
	ConversationID uuid.UUID  `gorm:"type:uuid"`
	UserID         *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
	Text           string
	Type           uint8
}

func (MessagePersistence) TableName() string {
	return "messages"
}

func ToMessageDTO(message *MessagePersistence, user *UserPersistence, requestUserID uuid.UUID) *MessageDTO {
	messageDTO := MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Text:      message.Text,
		Type:      "system",
	}

	if message.Type == 0 {
		messageDTO.User = ToUserDTO(user)
		messageDTO.UserId = message.UserID
		messageDTO.IsInbound = *message.UserID != requestUserID
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

func ToMessagePersistence(message *Message) *MessagePersistence {
	var messageType uint8

	switch message.Type {
	case "user":
		messageType = 0
	case "system":
		messageType = 1
	}

	return &MessagePersistence{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		CreatedAt:      message.CreatedAt,
		Text:           message.Text,
		Type:           messageType,
	}
}
