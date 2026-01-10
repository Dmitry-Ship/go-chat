package presentation

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	"github.com/google/uuid"
)

var messageTypesMap = map[uint8]domain.MessageType{
	0: domain.MessageTypeText,
	1: domain.MessageTypeRenamedConversation,
	2: domain.MessageTypeLeftConversation,
	3: domain.MessageTypeJoinedConversation,
	4: domain.MessageTypeInvitedConversation,
}

func MessageTypePersistenceToDomain(persistenceType uint8) domain.MessageType {
	if messageType, ok := messageTypesMap[persistenceType]; ok {
		return messageType
	}
	return domain.MessageTypeText
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
	case domain.MessageTypeText:
		return content
	case domain.MessageTypeRenamedConversation:
		return userName + " renamed chat to " + content
	case domain.MessageTypeLeftConversation:
		return userName + " left"
	case domain.MessageTypeJoinedConversation:
		return userName + " joined"
	case domain.MessageTypeInvitedConversation:
		return userName + " was invited"
	default:
		return content
	}
}

func (f *MessageFormatter) FormatMessageDTO(rawMessage readModel.RawMessageDTO, requestUserID uuid.UUID) readModel.MessageDTO {
	messageType := MessageTypePersistenceToDomain(rawMessage.Type)

	messageDTO := readModel.MessageDTO{
		ID:             rawMessage.ID,
		CreatedAt:      rawMessage.CreatedAt,
		Type:           messageType.String(),
		ConversationId: rawMessage.ConversationID,
		User: readModel.UserDTO{
			ID:     rawMessage.UserID,
			Avatar: rawMessage.UserAvatar,
			Name:   rawMessage.UserName,
		},
	}

	messageDTO.Text = f.FormatMessageText(messageType, rawMessage.Content, rawMessage.UserName)

	if messageType == domain.MessageTypeText {
		messageDTO.IsInbound = rawMessage.UserID != requestUserID
	}

	return messageDTO
}

func (f *MessageFormatter) FormatConversationLastMessage(rawLastMessage readModel.RawLastMessageDTO, userID uuid.UUID) readModel.MessageDTO {
	if rawLastMessage.MessageID == uuid.Nil {
		return readModel.MessageDTO{}
	}

	messageType := MessageTypePersistenceToDomain(uint8(rawLastMessage.MessageType))

	messageDTO := readModel.MessageDTO{
		ID:             rawLastMessage.MessageID,
		CreatedAt:      rawLastMessage.MessageCreatedAt,
		Type:           messageType.String(),
		ConversationId: rawLastMessage.ConversationID,
		User: readModel.UserDTO{
			ID:     rawLastMessage.MessageUserID,
			Avatar: rawLastMessage.MessageUserAvatar,
			Name:   rawLastMessage.MessageUserName,
		},
	}

	messageDTO.Text = f.FormatMessageText(messageType, rawLastMessage.MessageContent, rawLastMessage.MessageUserName)

	if messageType == domain.MessageTypeText {
		messageDTO.IsInbound = rawLastMessage.MessageUserID != userID
	}

	return messageDTO
}
