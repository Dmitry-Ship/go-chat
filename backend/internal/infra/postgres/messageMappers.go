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

func toMessagePersistence(message domain.BaseMessage) *Message {
	baseMessage := message.GetBaseData()
	return &Message{
		ID:             baseMessage.ID,
		ConversationID: baseMessage.ConversationID,
		UserID:         baseMessage.UserID,
		Type:           toMessageTypePersistence(baseMessage.Type),
	}
}

func toTextMessagePersistence(message domain.TextMessage) *TextMessage {
	return &TextMessage{
		ID:        message.ID,
		MessageID: message.GetBaseData().ID,
		Text:      message.Text,
	}
}

func toRenameConversationMessagePersistence(message domain.ConversationRenamedMessage) *ConversationRenamedMessage {
	return &ConversationRenamedMessage{
		ID:        message.ID,
		MessageID: message.GetBaseData().ID,
		NewName:   message.NewName,
	}
}
