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
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(toMessagePersistence(message)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(toTextMessagePersistence(*message)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	r.dispatchEvents(message)

	return nil
}

func (r *messageRepository) StoreLeftConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	r.dispatchEvents(message)

	return nil
}

func (r *messageRepository) StoreJoinedConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	r.dispatchEvents(message)

	return nil
}

func (r *messageRepository) StoreInvitedConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	r.dispatchEvents(message)

	return nil
}

func (r *messageRepository) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(toMessagePersistence(message)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(toRenameConversationMessagePersistence(*message)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	r.dispatchEvents(message)

	return nil
}
