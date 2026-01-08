package domain

import (
	"github.com/google/uuid"
)

type BaseConversation interface {
	GetBaseData() *Conversation
}

type ConversationType struct {
	slug string
}

func (r ConversationType) String() string {
	return r.slug
}

var (
	ConversationTypeGroup  = ConversationType{"group"}
	ConversationTypeDirect = ConversationType{"direct"}
)

type Conversation struct {
	ID       uuid.UUID
	Type     ConversationType
	IsActive bool
}

func (c *Conversation) GetBaseData() *Conversation {
	return c
}
