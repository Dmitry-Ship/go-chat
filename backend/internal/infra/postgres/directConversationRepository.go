package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type directConversationRepository struct {
	repository
}

func NewDirectConversationRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *directConversationRepository {
	return &directConversationRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *directConversationRepository) Store(conversation *domain.DirectConversation) error {
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

	if err := tx.Create(toDirectConversationPersistence(conversation)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(toParticipantPersistence(conversation.GetFromUser())).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(toParticipantPersistence(conversation.GetToUser())).Error; err != nil {
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

func (r *directConversationRepository) GetID(firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	directConversation := DirectConversation{}
	err := r.db.Where("to_user_id = ? AND from_user_id = ?", firstUserID, secondUserID).Or("to_user_id = ? AND from_user_id = ?", secondUserID, firstUserID).First(&directConversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return directConversation.ConversationID, nil
}
