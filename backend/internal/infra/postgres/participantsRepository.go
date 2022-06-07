package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantRepository struct {
	repository
}

func NewParticipantRepository(db *gorm.DB, eventPublisher infra.EventPublisher) *participantRepository {
	return &participantRepository{
		repository: *newRepository(db, eventPublisher),
	}
}

func (r *participantRepository) Store(participant *domain.Participant) error {
	return r.store(participant, toParticipantPersistence(*participant))
}

func (r *participantRepository) Update(participant *domain.Participant) error {
	return r.update(participant, toParticipantPersistence(*participant))
}

func (r *participantRepository) GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	var participantPersistence Participant

	err := r.db.Where(&Participant{ConversationID: conversationID, UserID: userID}).First(&participantPersistence).Error

	if err != nil {
		return nil, err
	}

	return toParticipantDomain(participantPersistence), nil
}

func (r *participantRepository) GetIDsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error) {
	var participants []Participant

	err := r.db.Where(&Participant{ConversationID: conversationID}).Find(&participants).Error

	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(participants))

	for i, participant := range participants {
		ids[i] = participant.UserID
	}

	return ids, nil
}
