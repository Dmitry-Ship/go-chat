package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
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
