package postgres

import (
	"GitHub/go-chat/backend/domain"

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
	err := r.participants.Create(domain.ToParticipantPersistence(participant)).Error

	return err
}

func (r *participantRepository) FindAllByConversationID(conversationID uuid.UUID) ([]*domain.Participant, error) {
	participants := []*domain.ParticipantPersistence{}

	err := r.participants.Limit(50).Where("conversation_id = ?", conversationID).Find(&participants).Error

	domainParticipants := make([]*domain.Participant, len(participants))

	for i, participant := range participants {
		domainParticipants[i] = domain.ToParticipantDomain(participant)
	}

	return domainParticipants, err
}

func (r *participantRepository) FindByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	participant := domain.ParticipantPersistence{}

	err := r.participants.Where("conversation_id = ?", conversationID).Where("user_id = ?", userID).First(&participant).Error

	return domain.ToParticipantDomain(&participant), err
}

func (r *participantRepository) DeleteByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) error {
	participant := domain.ParticipantPersistence{}

	err := r.participants.Where("conversation_id = ?", conversationID).Where("user_id = ?", userID).Delete(participant).Error

	return err
}
