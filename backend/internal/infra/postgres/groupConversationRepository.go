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
	return r.beginTransaction(conversation, func(tx *gorm.DB) error {
		if err := tx.Create(toConversationPersistence(conversation)).Error; err != nil {
			return err
		}

		if err := tx.Create(toGroupConversationPersistence(conversation)).Error; err != nil {
			return err
		}

		if err := tx.Create(toParticipantPersistence(&conversation.Owner)).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *groupConversationRepository) Update(conversation *domain.GroupConversation) error {
	return r.beginTransaction(conversation, func(tx *gorm.DB) error {
		if err := tx.Save(toConversationPersistence(conversation)).Error; err != nil {
			return err
		}

		if err := tx.Save(toGroupConversationPersistence(conversation)).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *groupConversationRepository) GetByID(id uuid.UUID) (*domain.GroupConversation, error) {
	conversation := Conversation{}

	err := r.db.Where(&Conversation{ID: id, IsActive: true}).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	groupConversation := GroupConversation{}

	err = r.db.Where(&GroupConversation{ConversationID: id}).First(&groupConversation).Error

	if err != nil {
		return nil, err
	}

	participant := Participant{}

	err = r.db.Where(&Participant{ConversationID: id, UserID: groupConversation.OwnerID}).First(&participant).Error

	if err != nil {
		return nil, err
	}

	return toGroupConversationDomain(&conversation, &groupConversation, &participant), nil
}
