package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationRepository struct {
	db             *gorm.DB
	eventPublisher infra.EventPublisher
}

func NewConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *conversationRepository {
	return &conversationRepository{
		db:             db,
		eventPublisher: eventPublisher,
	}

}

func (r *conversationRepository) StoreGroupConversation(conversation *domain.GroupConversation) error {
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

func (r *conversationRepository) UpdateGroupConversation(conversation *domain.GroupConversation) error {
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

func (r *conversationRepository) StoreDirectConversation(conversation *domain.DirectConversation) error {
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

func (r *conversationRepository) GetDirectConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	directConversation := DirectConversation{}
	err := r.db.Where("to_user_id = ? AND from_user_id = ?", firstUserId, secondUserID).Or("to_user_id = ? AND from_user_id = ?", secondUserID, firstUserId).First(&directConversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return directConversation.ConversationID, nil
}

func (r *conversationRepository) GetGroupConversation(id uuid.UUID) (*domain.GroupConversation, error) {
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
