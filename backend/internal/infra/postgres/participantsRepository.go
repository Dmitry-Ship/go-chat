package postgres

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *participantRepository {
	return &participantRepository{
		db: db,
	}
}

func (r *participantRepository) Store(participant *domain.Participant) error {
	err := r.db.Create(toParticipantPersistence(participant)).Error

	return err
}

func (r *participantRepository) DeleteByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) error {
	participant := Participant{}

	err := r.db.Where("conversation_id = ?", conversationID).Where("user_id = ?", userID).Delete(participant).Error

	return err
}
