package presentation

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	"github.com/google/uuid"
)

var messageTypesMap = map[uint8]domain.MessageType{
	0: domain.MessageTypeUser,
	1: domain.MessageTypeSystem,
}

func MessageTypePersistenceToDomain(persistenceType uint8) domain.MessageType {
	if messageType, ok := messageTypesMap[persistenceType]; ok {
		return messageType
	}
	return domain.MessageTypeUser
}

type MessageFormatter struct{}

func NewMessageFormatter() *MessageFormatter {
	return &MessageFormatter{}
}

func (f *MessageFormatter) FormatMessageText(messageType domain.MessageType, content string, userName string) string {
	if userName == "" {
		userName = "Unknown"
	}

	switch messageType {
	case domain.MessageTypeUser:
		return content
	case domain.MessageTypeSystem:
		if content != "" {
			return userName + " renamed chat to " + content
		}
		return userName + " performed a system action"
	default:
		return content
	}
}

func (f *MessageFormatter) FormatMessageDTO(rawMessage readModel.RawMessageDTO) readModel.MessageDTO {
	messageType := MessageTypePersistenceToDomain(rawMessage.Type)

	messageDTO := readModel.MessageDTO{
		ID:             rawMessage.ID,
		CreatedAt:      rawMessage.CreatedAt,
		Type:           messageType.String(),
		ConversationId: rawMessage.ConversationID,
		UserID:         rawMessage.UserID,
		User: &readModel.UserDTO{
			ID:     rawMessage.UserID,
			Avatar: rawMessage.UserAvatar,
			Name:   rawMessage.UserName,
		},
	}

	messageDTO.Text = f.FormatMessageText(messageType, rawMessage.Content, rawMessage.UserName)

	return messageDTO
}

func (f *MessageFormatter) FormatConversationLastMessage(rawLastMessage readModel.RawLastMessageDTO) readModel.MessageDTO {
	if rawLastMessage.MessageID == uuid.Nil {
		return readModel.MessageDTO{}
	}

	messageType := MessageTypePersistenceToDomain(uint8(rawLastMessage.MessageType))

	messageDTO := readModel.MessageDTO{
		ID:             rawLastMessage.MessageID,
		CreatedAt:      rawLastMessage.MessageCreatedAt,
		Type:           messageType.String(),
		ConversationId: rawLastMessage.ConversationID,
		UserID:         rawLastMessage.MessageUserID,
		User: &readModel.UserDTO{
			ID:     rawLastMessage.MessageUserID,
			Avatar: rawLastMessage.MessageUserAvatar,
			Name:   rawLastMessage.MessageUserName,
		},
	}

	messageDTO.Text = f.FormatMessageText(messageType, rawLastMessage.MessageContent, rawLastMessage.MessageUserName)

	return messageDTO
}
