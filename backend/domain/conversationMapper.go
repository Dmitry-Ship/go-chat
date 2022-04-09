package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationDTOFull struct {
	Conversation ConversationDTO `json:"conversation"`
	HasJoined    bool            `json:"joined"`
}

type ConversationPersistence struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	Name      string
	Avatar    string
	IsPrivate bool
	CreatedAt time.Time
}

func (ConversationPersistence) TableName() string {
	return "conversations"
}

func ToConversationDTOFull(conversation *ConversationPersistence, hasJoined bool) *ConversationDTOFull {
	return &ConversationDTOFull{
		Conversation: *ToConversationDTO(conversation),
		HasJoined:    hasJoined,
	}
}

func ToConversationDTO(conversation *ConversationPersistence) *ConversationDTO {
	return &ConversationDTO{
		ID:        conversation.ID,
		Name:      conversation.Name,
		Avatar:    conversation.Avatar,
		IsPrivate: conversation.IsPrivate,
		CreatedAt: conversation.CreatedAt,
	}
}

func ToConversationPersistence(conversation *Conversation) *ConversationPersistence {
	return &ConversationPersistence{
		ID:        conversation.ID,
		Name:      conversation.Name,
		Avatar:    conversation.Avatar,
		IsPrivate: conversation.IsPrivate,
		CreatedAt: conversation.CreatedAt,
	}
}
