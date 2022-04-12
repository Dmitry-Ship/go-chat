package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
	Type           uint8
}

type TextMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	Text      string
}

type ConversationRenamedMessage struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	NewName   string
}

func toMessageDTO(message *Message, user *User) *readModel.MessageDTO {
	var messageType string
	switch message.Type {
	case 0:
		messageType = "text"
	case 1:
		messageType = "renamed_conversation"
	case 2:
		messageType = "left_conversation"
	case 3:
		messageType = "joined_conversation"
	}

	messageDTO := readModel.MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Type:      messageType,
		User:      toUserDTO(user),
	}

	return &messageDTO
}

func ToTextMessageDTO(message *Message, user *User, text string, requestUserID uuid.UUID) *readModel.MessageDTO {
	messageDTO := toMessageDTO(message, user)
	messageDTO.IsInbound = user.ID != requestUserID
	messageDTO.Text = text

	return messageDTO
}

func toConversationRenamedMessageDTO(message *Message, user *User, newName string) *readModel.MessageDTO {
	messageDTO := toMessageDTO(message, user)
	messageDTO.NewConversationName = newName

	return messageDTO
}

func toMessagePersistence(message domain.BaseMessage) *Message {
	var messageType uint8

	baseMessage := message.GetBaseData()

	switch baseMessage.Type {
	case "text":
		messageType = 0
	case "renamed_conversation":
		messageType = 1
	case "left_conversation":
		messageType = 2
	case "joined_conversation":
		messageType = 3
	}

	return &Message{
		ID:             baseMessage.ID,
		ConversationID: baseMessage.ConversationID,
		UserID:         baseMessage.UserID,
		CreatedAt:      baseMessage.CreatedAt,
		Type:           messageType,
	}
}

func toTextMessagePersistence(message domain.TextMessage) *TextMessage {
	text := message.GetTextMessageData()

	return &TextMessage{
		ID:        text.ID,
		MessageID: text.MessageID,
		Text:      text.Text,
	}
}

func toRenameConversationMessagePersistence(message domain.ConversationRenamedMessage) *ConversationRenamedMessage {
	conversationRenamedMessage := message.GetConversationRenamedMessage()

	return &ConversationRenamedMessage{
		ID:        conversationRenamedMessage.ID,
		MessageID: conversationRenamedMessage.MessageID,
		NewName:   conversationRenamedMessage.NewName,
	}
}
