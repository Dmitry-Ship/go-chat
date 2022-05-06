package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"gorm.io/gorm"
)

type messageRepository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func NewMessageRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *messageRepository {
	return &messageRepository{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (r *messageRepository) StoreTextMessage(message *domain.TextMessage) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toTextMessagePersistence(*message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.eventPublisher)

	return nil
}

func (r *messageRepository) StoreLeftConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.eventPublisher)

	return nil
}

func (r *messageRepository) StoreJoinedConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.eventPublisher)

	return nil
}

func (r *messageRepository) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toRenameConversationMessagePersistence(*message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.eventPublisher)

	return nil
}
