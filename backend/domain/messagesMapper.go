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

func ToMessageDTO(message *MessagePersistence) *MessageDTO {
	var messageType string

	switch message.Type {
	case 0:
		messageType = "user"
	case 1:
		messageType = "system"
	}

	return &MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Text:      message.Text,
		Type:      messageType,
	}
}

func ToMessageDTOFromDOmain(message *Message) *MessageDTO {
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
