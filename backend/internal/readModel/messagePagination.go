package readModel

import (
	"time"

	"github.com/google/uuid"
)

type MessageCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

type MessagePageDTO struct {
	Messages   []MessageDTO `json:"messages"`
	NextCursor string       `json:"next_cursor,omitempty"`
	HasMore    bool         `json:"has_more"`
}
