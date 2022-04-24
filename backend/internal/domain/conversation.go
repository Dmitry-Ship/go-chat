package domain

import (
	"time"

	"github.com/google/uuid"
)

type BaseConversation interface {
	GetBaseData() *Conversation
}

type Conversation struct {
	ID        uuid.UUID
	Type      string
	CreatedAt time.Time
}

func (c *Conversation) GetBaseData() *Conversation {
	return c
}

type PublicConversationData struct {
	ID     uuid.UUID
	Name   string
	Avatar string
}
type PublicConversation struct {
	Conversation
	Data PublicConversationData
}

func NewPublicConversation(id uuid.UUID, name string) *PublicConversation {
	return &PublicConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "public",
			CreatedAt: time.Now(),
		},
		Data: PublicConversationData{
			ID:     uuid.New(),
			Name:   name,
			Avatar: string(name[0]),
		},
	}
}

func (publicConversation *PublicConversation) Rename(newName string) {
	publicConversation.Data.Name = newName
}

type PrivateConversationData struct {
	ID         uuid.UUID
	ToUserId   uuid.UUID
	FromUserId uuid.UUID
}

type PrivateConversation struct {
	Conversation
	Data PrivateConversationData
}

func NewPrivateConversation(id uuid.UUID, to uuid.UUID, from uuid.UUID) *PrivateConversation {
	return &PrivateConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "private",
			CreatedAt: time.Now(),
		},
		Data: PrivateConversationData{
			ID:         uuid.New(),
			ToUserId:   to,
			FromUserId: from,
		},
	}
}
