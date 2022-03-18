package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid" json:"id"`
	RoomId    uuid.UUID `json:"room_id"`
	Content   string    `json:"content"`
	CreatedAt int32     `json:"created_at"`
	Type      string    `json:"type"`
	UserId    uuid.UUID `json:"user_id"`
}

func NewChatMessage(content string, messageType string, roomId uuid.UUID, userId uuid.UUID) *ChatMessage {
	return &ChatMessage{
		ID:        uuid.New(),
		RoomId:    roomId,
		Content:   content,
		CreatedAt: int32(time.Now().Unix()),
		Type:      messageType,
		UserId:    userId,
	}
}
