package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantRepository struct {
	participants *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *participantRepository {
	return &participantRepository{
		participants: db,
	}
}

func (r *participantRepository) Store(participant *domain.Participant) error {
	err := r.participants.Create(mappers.ToParticipantPersistence(participant)).Error

	return err
}

func (r *participantRepository) FindAllByConversationID(conversationID uuid.UUID) ([]*domain.Participant, error) {
	participants := []*mappers.ParticipantPersistence{}

	err := r.participants.Limit(50).Where("conversation_id = ?", conversationID).Find(&participants).Error

	domainParticipants := make([]*domain.Participant, len(participants))

	for i, participant := range participants {
		domainParticipants[i] = mappers.ToParticipantDomain(participant)
	}

	return domainParticipants, err
}

func (r *participantRepository) FindByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	participant := mappers.ParticipantPersistence{}

	err := r.participants.Where("conversation_id = ?", conversationID).Where("user_id = ?", userID).First(&participant).Error

	return mappers.ToParticipantDomain(&participant), err
}

func (r *participantRepository) DeleteByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) error {
	participant := mappers.ParticipantPersistence{}

	err := r.participants.Where("conversation_id = ?", conversationID).Where("user_id = ?", userID).Delete(participant).Error

	return err
}
