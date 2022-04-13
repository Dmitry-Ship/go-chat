package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	Name      string
	Avatar    string
	Type      uint8
	CreatedAt time.Time
}

var conversationTypesMap = map[uint8]string{
	0: "public",
	1: "private",
}

func toConversationTypePersistence(conversationType string) uint8 {
	for k, v := range conversationTypesMap {
		if v == conversationType {
			return k
		}
	}

	return 0
}

func toConversationDTO(conversation *Conversation) *readModel.ConversationDTO {
	return &readModel.ConversationDTO{
		ID:        conversation.ID,
		Name:      conversation.Name,
		Avatar:    conversation.Avatar,
		Type:      conversationTypesMap[conversation.Type],
		CreatedAt: conversation.CreatedAt,
	}
}

func toConversationFullDTO(conversation *Conversation, hasJoined bool) *readModel.ConversationFullDTO {
	return &readModel.ConversationFullDTO{
		Conversation: *toConversationDTO(conversation),
		HasJoined:    hasJoined,
	}
}

func toConversationPersistence(conversation *domain.Conversation) *Conversation {
	return &Conversation{
		ID:        conversation.ID,
		Name:      conversation.Name,
		Avatar:    conversation.Avatar,
		Type:      toConversationTypePersistence(conversation.Type),
		CreatedAt: conversation.CreatedAt,
	}
}
