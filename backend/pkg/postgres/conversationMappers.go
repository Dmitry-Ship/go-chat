package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID uuid.UUID `gorm:"type:uuid"`

	Type      uint8
	CreatedAt time.Time
}

type PublicConversation struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	Name           string
	Avatar         string
	ConversationID uuid.UUID `gorm:"type:uuid"`
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
	dto := readModel.ConversationDTO{
		ID:        conversation.ID,
		CreatedAt: conversation.CreatedAt,
		Type:      messageTypesMap[conversation.Type],
	}

	return &dto
}

func toPublicConversationDTO(conversation *Conversation, avatar string, name string) *readModel.ConversationDTO {
	dto := toConversationDTO(conversation)

	dto.Avatar = avatar
	dto.Name = name

	return dto
}

func toConversationFullDTO(conversation *Conversation, avatar string, name string, hasJoined bool) *readModel.ConversationFullDTO {
	return &readModel.ConversationFullDTO{
		Conversation: *toPublicConversationDTO(conversation, avatar, name),
		HasJoined:    hasJoined,
	}
}

func toConversationPersistence(conversation domain.BaseConversation) *Conversation {
	conversationBase := conversation.GetBaseData()
	return &Conversation{
		ID:        conversationBase.ID,
		Type:      toConversationTypePersistence(conversationBase.Type),
		CreatedAt: conversationBase.CreatedAt,
	}
}

func toPublicConversationPersistence(conversation *domain.PublicConversation) *PublicConversation {
	return &PublicConversation{
		ID:             conversation.Data.ID,
		ConversationID: conversation.ID,
		Name:           conversation.Data.Name,
		Avatar:         conversation.Data.Avatar,
	}
}
