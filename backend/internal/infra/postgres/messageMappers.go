package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	"time"

	"github.com/google/uuid"
)

var messageTypesMap = map[uint8]domain.MessageType{
	0: domain.MessageTypeText,
	1: domain.MessageTypeRenamedConversation,
	2: domain.MessageTypeLeftConversation,
	3: domain.MessageTypeJoinedConversation,
	4: domain.MessageTypeInvitedConversation,
}

func toMessageTypePersistence(messageType domain.MessageType) uint8 {
	for k, v := range messageTypesMap {
		if v == messageType {
			return k
		}
	}

	return 0
}

func toMessagePersistence(message *domain.Message) *Message {
	return &Message{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		Type:           toMessageTypePersistence(message.Type),
		Content:        message.Content.String(),
	}
}

type messageQuery struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	Type           uint8
	UserID         uuid.UUID
	ConversationID uuid.UUID
	UserName       string
	UserAvatar     string
	Content        string
}

func toMessageDTO(message messageQuery, requestUserID uuid.UUID) readModel.MessageDTO {
	text := ""

	switch messageTypesMap[message.Type] {
	case domain.MessageTypeText:
		text = message.Content
	case domain.MessageTypeRenamedConversation:
		text = message.UserName + " renamed chat to " + message.Content
	case domain.MessageTypeJoinedConversation:
		text = message.UserName + " joined"
	case domain.MessageTypeLeftConversation:
		text = message.UserName + " left"
	case domain.MessageTypeInvitedConversation:
		text = message.UserName + " was invited"
	default:
		text = "Unknown message type"
	}

	massageDTO := readModel.MessageDTO{
		ID:             message.ID,
		CreatedAt:      message.CreatedAt,
		Text:           text,
		Type:           messageTypesMap[message.Type].String(),
		ConversationId: message.ConversationID,
		User: readModel.UserDTO{
			ID:     message.UserID,
			Avatar: message.UserAvatar,
			Name:   message.UserName,
		},
	}

	if messageTypesMap[message.Type] == domain.MessageTypeText {
		massageDTO.IsInbound = message.UserID != requestUserID
	}

	return massageDTO
}
