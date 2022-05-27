package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

var messageTypesMap = map[uint8]string{
	0: domain.MessageTypeText,
	1: domain.MessageTypeRenamedConversation,
	2: domain.MessageTypeLeftConversation,
	3: domain.MessageTypeJoinedConversation,
	4: domain.MessageTypeInvitedConversation,
}

func toMessageTypePersistence(messageType string) uint8 {
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
