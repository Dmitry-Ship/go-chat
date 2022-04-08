package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"

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

	err := r.chatMessages.Create(mappers.ToMessagePersistence(chatMessage)).Error

	return err
}

func (r *messageRepository) FindAllByConversationID(conversationID uuid.UUID) ([]*domain.Message, error) {
	messages := []*mappers.MessagePersistence{}

	err := r.chatMessages.Limit(50).Where("conversation_id = ?", conversationID).Find(&messages).Error

	domainMessages := make([]*domain.Message, len(messages))

	for i, message := range messages {
		domainMessages[i] = mappers.ToMessageDomain(message)
	}

	return domainMessages, err
}
