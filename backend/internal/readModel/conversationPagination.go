package readModel

import (
	"time"

	"github.com/google/uuid"
)

type ConversationCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

type ConversationPageDTO struct {
	Conversations []ConversationDTO `json:"conversations"`
	NextCursor    string            `json:"next_cursor,omitempty"`
	HasMore       bool              `json:"has_more"`
}
