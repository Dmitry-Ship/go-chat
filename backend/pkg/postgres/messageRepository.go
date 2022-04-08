package postgres

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	chatMessages *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{
		chatMessages: db,
	}
}

func (r *messageRepository) Store(chatMessage *domain.Message) error {

	err := r.chatMessages.Create(domain.ToMessagePersistence(chatMessage)).Error

	return err
}

func (r *messageRepository) FindAllByConversationID(conversationID uuid.UUID) ([]*domain.MessageDTO, error) {
	messages := []*domain.MessagePersistence{}

	err := r.chatMessages.Limit(50).Where("conversation_id = ?", conversationID).Find(&messages).Error

	dtoMessages := make([]*domain.MessageDTO, len(messages))

	for i, message := range messages {
		dtoMessages[i] = domain.ToMessageDTO(message)
	}

	return dtoMessages, err
}
