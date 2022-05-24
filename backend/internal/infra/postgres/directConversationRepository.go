package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type directConversationRepository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func NewDirectConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *directConversationRepository {
	return &directConversationRepository{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (r *directConversationRepository) GetByID(id uuid.UUID) (*domain.DirectConversation, error) {
	conversation := Conversation{}

	err := r.db.Where("id = ?", id).Where("is_active = ?", true).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	directConversation := DirectConversation{}

	err = r.db.Where("conversation_id = ?", id).First(&directConversation).Error

	if err != nil {
		return nil, err
	}

	toUser := Participant{}

	err = r.db.Where("conversation_id = ? AND user_id = ?", id, directConversation.ToUserID).First(&toUser).Error

	if err != nil {
		return nil, err
	}

	fromUser := Participant{}

	err = r.db.Where("conversation_id = ? AND user_id = ?", id, directConversation.FromUserID).First(&fromUser).Error

	if err != nil {
		return nil, err
	}

	return toDirectConversationDomain(&conversation, &directConversation, &toUser, &fromUser), nil
}

func (r *directConversationRepository) Store(conversation *domain.DirectConversation) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toDirectConversationPersistence(conversation)).Error

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

func (r *directConversationRepository) GetID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	directConversation := DirectConversation{}
	err := r.db.Where("to_user_id = ? AND from_user_id = ?", firstUserId, secondUserID).Or("to_user_id = ? AND from_user_id = ?", secondUserID, firstUserId).First(&directConversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return directConversation.ConversationID, nil
}
