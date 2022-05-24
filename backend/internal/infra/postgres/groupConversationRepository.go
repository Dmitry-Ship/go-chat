package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type groupConversationRepository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func NewGroupConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *groupConversationRepository {
	return &groupConversationRepository{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (r *groupConversationRepository) Store(conversation *domain.GroupConversation) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toGroupConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toParticipantPersistence(&conversation.Data.Owner)).Error

	if err != nil {
		return err
	}

	conversation.Dispatch(r.eventPublisher)

	return nil
}

func (r *groupConversationRepository) Update(conversation *domain.GroupConversation) error {
	err := r.db.Save(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Save(toGroupConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	conversation.Dispatch(r.eventPublisher)

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
