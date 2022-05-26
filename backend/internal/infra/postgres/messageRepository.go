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

func (r *messageRepository) StoreTextMessage(message *domain.TextMessage) error {
	return r.beginTransaction(message, func(tx *gorm.DB) error {
		if err := tx.Create(toMessagePersistence(message)).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Create(toTextMessagePersistence(*message)).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func (r *messageRepository) StoreLeftConversationMessage(message *domain.Message) error {
	return r.store(message, toMessagePersistence(message))
}

func (r *messageRepository) StoreJoinedConversationMessage(message *domain.Message) error {
	return r.store(message, toMessagePersistence(message))
}

func (r *messageRepository) StoreInvitedConversationMessage(message *domain.Message) error {
	return r.store(message, toMessagePersistence(message))
}

func (r *messageRepository) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	return r.beginTransaction(message, func(tx *gorm.DB) error {
		if err := tx.Create(toMessagePersistence(message)).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Create(toRenameConversationMessagePersistence(*message)).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}
