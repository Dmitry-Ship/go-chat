package postgres

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantRepository struct {
	db     *gorm.DB
	pubsub domain.EventPublisher
}

func NewParticipantRepository(db *gorm.DB, pubsub domain.EventPublisher) *participantRepository {
	return &participantRepository{
		db:     db,
		pubsub: pubsub,
	}
}

func (r *participantRepository) Store(participant *domain.Participant) error {
	err := r.db.Create(toParticipantPersistence(participant)).Error

	if err != nil {
		return err
	}

	participant.Dispatch(r.pubsub)

	return nil
}

func (r *participantRepository) Update(participant *domain.Participant) error {
	err := r.db.Save(toParticipantPersistence(participant)).Error

	if err != nil {
		return err
	}

	participant.Dispatch(r.pubsub)

	return nil
}

func (r *participantRepository) GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	var participantPersistence Participant

	err := r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&participantPersistence).Error

	if err != nil {
		return nil, err
	}

	return toParticipantDomain(&participantPersistence), nil
}
