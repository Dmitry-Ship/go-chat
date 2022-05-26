package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type groupConversationRepository struct {
	repository
}

func NewGroupConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *groupConversationRepository {
	return &groupConversationRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *groupConversationRepository) Store(conversation *domain.GroupConversation) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(toConversationPersistence(conversation)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(toGroupConversationPersistence(conversation)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(toParticipantPersistence(&conversation.Data.Owner)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	r.dispatchEvents(conversation)

	return nil
}

func (r *groupConversationRepository) Update(conversation *domain.GroupConversation) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Save(toConversationPersistence(conversation)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(toGroupConversationPersistence(conversation)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	r.dispatchEvents(conversation)

	return nil
}

func (r *groupConversationRepository) GetByID(id uuid.UUID) (*domain.GroupConversation, error) {
	conversation := Conversation{}

	err := r.db.Where("id = ?", id).Where("is_active = ?", true).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	groupConversation := GroupConversation{}

	err = r.db.Where("conversation_id = ?", id).First(&groupConversation).Error

	if err != nil {
		return nil, err
	}

	participant := Participant{}

	err = r.db.Where("conversation_id = ? AND user_id = ?", id, groupConversation.OwnerID).First(&participant).Error

	if err != nil {
		return nil, err
	}

	return toGroupConversationDomain(&conversation, &groupConversation, &participant), nil
}
