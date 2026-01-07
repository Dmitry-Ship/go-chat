package postgres

import (
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type directConversationRepository struct {
	repository
}

func NewDirectConversationRepository(db *gorm.DB, eventPublisher *infra.EventBus) *directConversationRepository {
	return &directConversationRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *directConversationRepository) Store(conversation *domain.DirectConversation) error {
	return r.beginTransaction(conversation, func(tx *gorm.DB) error {
		if err := tx.Create(toConversationPersistence(conversation)).Error; err != nil {
			return fmt.Errorf("create conversation error: %w", err)
		}

		for _, participant := range conversation.Participants {
			if err := tx.Create(&Participant{
				ID:             participant.ID,
				ConversationID: participant.ConversationID,
				UserID:         participant.UserID,
				IsActive:       participant.IsActive,
			}).Error; err != nil {
				return fmt.Errorf("create participant error: %w", err)
			}
		}

		return nil
	})
}

func (r *directConversationRepository) GetByID(id uuid.UUID) (*domain.DirectConversation, error) {
	type result struct {
		Conversation  Conversation
		ParticipantID uuid.UUID
		UserID        uuid.UUID
	}

	results := []*result{}

	err := r.db.
		Table("conversations").
		Select("conversations.*, participants.id as participant_id, participants.user_id").
		Joins("LEFT JOIN participants ON participants.conversation_id = conversations.id").
		Where("conversations.id = ? AND conversations.is_active = ?", id, true).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("get direct conversation error: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("direct conversation not found")
	}

	participants := make([]*Participant, len(results))
	for i, res := range results {
		participants[i] = &Participant{
			ID:             res.ParticipantID,
			ConversationID: results[0].Conversation.ID,
			UserID:         res.UserID,
		}
	}

	return toDirectConversationDomain(&results[0].Conversation, participants), nil
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
		return uuid.Nil, fmt.Errorf("get conversation id error: %w", err)
	}

	return conversation.ID, nil
}
