package postgres

import (
	"fmt"

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
			return fmt.Errorf("create conversation error: %w", err)
		}

		if err := tx.Create(toGroupConversationPersistence(conversation)).Error; err != nil {
			return fmt.Errorf("create group conversation error: %w", err)
		}

		if err := tx.Create(toParticipantPersistence(conversation.Owner)).Error; err != nil {
			return fmt.Errorf("create participant error: %w", err)
		}

		return nil
	})
}

func (r *groupConversationRepository) Update(conversation *domain.GroupConversation) error {
	return r.beginTransaction(conversation, func(tx *gorm.DB) error {
		if err := tx.Save(toConversationPersistence(conversation)).Error; err != nil {
			return fmt.Errorf("save conversation error: %w", err)
		}

		if err := tx.Save(toGroupConversationPersistence(conversation)).Error; err != nil {
			return fmt.Errorf("save group conversation error: %w", err)
		}

		return nil
	})
}

func (r *groupConversationRepository) GetByID(id uuid.UUID) (*domain.GroupConversation, error) {
	type result struct {
		Conversation      Conversation
		GroupConversation GroupConversation
		Participant       Participant
	}

	var res result

	err := r.db.
		Table("conversations").
		Select("conversations.*, group_conversations.*, participants.*").
		Joins("INNER JOIN group_conversations ON group_conversations.conversation_id = conversations.id").
		Joins("INNER JOIN participants ON participants.conversation_id = conversations.id AND participants.user_id = group_conversations.owner_id").
		Where("conversations.id = ? AND conversations.is_active = ?", id, true).
		Scan(&res).Error

	if err != nil {
		return nil, fmt.Errorf("get group conversation error: %w", err)
	}

	return toGroupConversationDomain(res.Conversation, res.GroupConversation, res.Participant), nil
}
