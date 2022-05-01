package postgres

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationRepository struct {
	db     *gorm.DB
	pubsub domain.EventPublisher
}

func NewConversationRepository(db *gorm.DB, pubsub domain.EventPublisher) *conversationRepository {
	return &conversationRepository{
		db:     db,
		pubsub: pubsub,
	}

}

func (r *conversationRepository) StorePublicConversation(conversation *domain.PublicConversation) error {
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

	conversation.Dispatch(r.pubsub)

	return nil
}

func (r *conversationRepository) UpdatePublicConversation(conversation *domain.PublicConversation) error {
	err := r.db.Save(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Save(toPublicConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	conversation.Dispatch(r.pubsub)

	return nil
}

func (r *conversationRepository) StorePrivateConversation(conversation *domain.PrivateConversation) error {
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

	conversation.Dispatch(r.pubsub)

	return nil
}

func (r *conversationRepository) GetPrivateConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	privateConversation := PrivateConversation{}
	err := r.db.Where("to_user_id = ? AND from_user_id = ?", firstUserId, secondUserID).Or("to_user_id = ? AND from_user_id = ?", secondUserID, firstUserId).First(&privateConversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return privateConversation.ConversationID, nil
}

func (r *conversationRepository) GetPublicConversation(id uuid.UUID) (*domain.PublicConversation, error) {
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
