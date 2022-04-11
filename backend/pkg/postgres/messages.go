package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID  `gorm:"type:uuid"`
	ConversationID uuid.UUID  `gorm:"type:uuid"`
	UserID         *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
	Text           string
	Type           uint8
}

type TextMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	Text      string
}

type LeftMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
}

type JoinedMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
}

type ConversationRenamedMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	NewName   string
}

func ToMessageDTO(message *Message, user *User, requestUserID uuid.UUID) *readModel.MessageDTO {
	messageDTO := readModel.MessageDTO{
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

func ToMessagePersistence(message *domain.Message) *Message {
	var messageType uint8

	switch message.Type {
	case "user":
		messageType = 0
	case "system":
		messageType = 1
	}

	return &Message{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		CreatedAt:      message.CreatedAt,
		Text:           message.Text,
		Type:           messageType,
	}
}
