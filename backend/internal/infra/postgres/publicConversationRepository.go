package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type publicConversationRepository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func NewPublicConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *publicConversationRepository {
	return &publicConversationRepository{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (r *publicConversationRepository) Store(conversation *domain.PublicConversation) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toPublicConversationPersistence(conversation)).Error

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

func (r *publicConversationRepository) Update(conversation *domain.PublicConversation) error {
	err := r.db.Save(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Save(toPublicConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	conversation.Dispatch(r.eventPublisher)

	return nil
}

func (r *publicConversationRepository) GetByID(id uuid.UUID) (*domain.PublicConversation, error) {
	conversation := Conversation{}

	err := r.db.Where("id = ?", id).Where("is_active = ?", true).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	publicConversation := PublicConversation{}

	err = r.db.Where("conversation_id = ?", id).First(&publicConversation).Error

	if err != nil {
		return nil, err
	}

	participant := Participant{}

	err = r.db.Where("conversation_id = ? AND user_id = ?", id, publicConversation.OwnerID).First(&participant).Error

	if err != nil {
		return nil, err
	}

	return toPublicConversationDomain(&conversation, &publicConversation, &participant), nil
}
