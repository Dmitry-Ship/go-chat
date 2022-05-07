package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type privateConversationRepository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func NewPrivateConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *privateConversationRepository {
	return &privateConversationRepository{
		db:             db,
		eventPublisher: eventPublisher,
	}

}

func (r *privateConversationRepository) Store(conversation *domain.PrivateConversation) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toPrivateConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toParticipantPersistence(conversation.GetFromUser())).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toParticipantPersistence(conversation.GetToUser())).Error

	if err != nil {
		return err
	}

	conversation.Dispatch(r.eventPublisher)

	return nil
}

func (r *privateConversationRepository) GetID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	privateConversation := PrivateConversation{}
	err := r.db.Where("to_user_id = ? AND from_user_id = ?", firstUserId, secondUserID).Or("to_user_id = ? AND from_user_id = ?", secondUserID, firstUserId).First(&privateConversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return privateConversation.ConversationID, nil
}
