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
	IsPrivate bool
	CreatedAt time.Time
}

func ToConversationDTOFull(conversation *Conversation, hasJoined bool) *readModel.ConversationDTOFull {
	return &readModel.ConversationDTOFull{
		Conversation: *ToConversationDTO(conversation),
		HasJoined:    hasJoined,
	}
}

func ToConversationDTO(conversation *Conversation) *readModel.ConversationDTO {
	return &readModel.ConversationDTO{
		ID:        conversation.ID,
		Name:      conversation.Name,
		Avatar:    conversation.Avatar,
		IsPrivate: conversation.IsPrivate,
		CreatedAt: conversation.CreatedAt,
	}
}

func ToConversationPersistence(conversation *domain.Conversation) *Conversation {
	return &Conversation{
		ID:        conversation.ID,
		Name:      conversation.Name,
		Avatar:    conversation.Avatar,
		IsPrivate: conversation.IsPrivate,
		CreatedAt: conversation.CreatedAt,
	}
}
