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

var messageTypesMap = map[uint8]string{
	0: "text",
	1: "conversation_renamed",
	2: "left_conversation",
	3: "joined_conversation",
}

func toMessageTypePersistence(messageType string) uint8 {
	for k, v := range messageTypesMap {
		if v == messageType {
			return k
		}
	}

	return 0
}

func toMessageDTO(message *Message, user *User) *readModel.MessageDTO {
	messageDTO := readModel.MessageDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		Type:      messageTypesMap[message.Type],
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
	baseMessage := message.GetBaseData()
	return &Message{
		ID:             baseMessage.ID,
		ConversationID: baseMessage.ConversationID,
		UserID:         baseMessage.UserID,
		CreatedAt:      baseMessage.CreatedAt,
		Type:           toMessageTypePersistence(baseMessage.Type),
	}
}

func toTextMessagePersistence(message domain.TextMessage) *TextMessage {
	text := message.GetTextMessageData()

	return &TextMessage{
		ID:        text.ID,
		MessageID: message.GetBaseData().ID,
		Text:      text.Text,
	}
}

func toRenameConversationMessagePersistence(message domain.ConversationRenamedMessage) *ConversationRenamedMessage {
	conversationRenamedMessage := message.GetConversationRenamedMessage()

	return &ConversationRenamedMessage{
		ID:        conversationRenamedMessage.ID,
		MessageID: message.GetBaseData().ID,
		NewName:   conversationRenamedMessage.NewName,
	}
}
