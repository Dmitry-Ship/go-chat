package postgres

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type chatMessageRepository struct {
	chatMessages *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) *chatMessageRepository {
	return &chatMessageRepository{
		chatMessages: db,
	}
}

func (r *chatMessageRepository) Store(chatMessage *domain.ChatMessage) error {
	err := r.chatMessages.Create(&chatMessage).Error

	return err
}

func (r *chatMessageRepository) FindAllByRoomID(roomID uuid.UUID) ([]*domain.ChatMessage, error) {
	messages := []*domain.ChatMessage{}

	err := r.chatMessages.Limit(50).Where("room_id = ?", roomID).Find(&messages).Error

	return messages, err
}
