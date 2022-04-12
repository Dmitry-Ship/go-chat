package postgres

import (
	"GitHub/go-chat/backend/domain"

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

func (r *participantRepository) Store(participant *domain.ParticipantAggregate) error {
	err := r.db.Create(ToParticipantPersistence(participant)).Error

	return err
}

func (r *participantRepository) GetUserIdsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.db.Model(&Participant{}).Where("conversation_id = ?", conversationID).Select("user_id").Find(&ids).Error

	return ids, err
}

func (r *participantRepository) DeleteByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) error {
	participant := Participant{}

	err := r.db.Where("conversation_id = ?", conversationID).Where("user_id = ?", userID).Delete(participant).Error

	return err
}
