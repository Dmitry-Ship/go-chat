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
	return r.beginTransaction(conversation, func(tx *gorm.DB) error {
		if err := tx.Create(toConversationPersistence(conversation)).Error; err != nil {
			return err
		}

		for _, participant := range conversation.Participants {
			if err := tx.Create(toParticipantPersistence(participant)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *directConversationRepository) GetByID(id uuid.UUID) (*domain.DirectConversation, error) {
	conversation := Conversation{}

	err := r.db.Where(&Conversation{ID: id, IsActive: true}).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	participants := []*Participant{}

	err = r.db.Where(&Participant{ConversationID: id}).Find(&participants).Error

	if err != nil {
		return nil, err
	}

	return toDirectConversationDomain(&conversation, participants), nil
}

func (r *directConversationRepository) GetID(firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	conversation := Conversation{}

	err := r.db.
		Model(&Conversation{}).
		Joins("INNER JOIN participants ON participants.conversation_id = conversations.id").
		Where(&Conversation{IsActive: true, Type: toConversationTypePersistence(domain.ConversationTypeDirect)}).
		Where("participants.is_active = ?", true).
		Where(r.db.Where("participants.user_id = ? ", secondUserID).Or("participants.user_id = ? ", secondUserID)).
		First(&conversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return conversation.ID, nil
}
