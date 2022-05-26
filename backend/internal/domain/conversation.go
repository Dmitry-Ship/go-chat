package domain

import (
	"time"

	"github.com/google/uuid"
)

type BaseConversation interface {
	GetBaseData() *Conversation
}
type Conversation struct {
	aggregate
	ID        uuid.UUID
	Type      string
	CreatedAt time.Time
	IsActive  bool
}

func (c *Conversation) GetBaseData() *Conversation {
	return c
}
