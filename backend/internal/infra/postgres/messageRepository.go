package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"gorm.io/gorm"
)

type messageRepository struct {
	repository
}

func NewMessageRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *messageRepository {
	return &messageRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *messageRepository) Store(message *domain.Message) error {
	return r.repository.store(message, toMessagePersistence(message))
}
