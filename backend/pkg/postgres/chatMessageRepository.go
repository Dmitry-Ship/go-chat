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

func (r *chatMessageRepository) Store(chatMessage *domain.Message) error {
	err := r.chatMessages.Create(&chatMessage).Error

	return err
}

func (r *chatMessageRepository) FindAllByConversationID(conversationID uuid.UUID) ([]*domain.Message, error) {
	messages := []*domain.Message{}

	err := r.chatMessages.Limit(50).Where("conversation_id = ?", conversationID).Find(&messages).Error

	return messages, err
}
