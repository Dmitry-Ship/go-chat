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
	err := r.participants.Create(&participant).Error

	return err
}

func (r *participantRepository) FindAllByRoomID(roomID uuid.UUID) ([]*domain.Participant, error) {
	participants := []*domain.Participant{}

	err := r.participants.Limit(50).Where("room_id = ?", roomID).Find(&participants).Error

	return participants, err
}

func (r *participantRepository) FindByRoomIDAndUserID(roomID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	participant := domain.Participant{}
	err := r.participants.Where("room_id = ?", roomID).Where("user_id = ?", userID).First(&participant).Error

	return &participant, err
}

func (r *participantRepository) DeleteByRoomIDAndUserID(roomID uuid.UUID, userID uuid.UUID) error {
	participant := domain.Participant{}

	err := r.participants.Where("room_id = ?", roomID).Where("user_id = ?", userID).Delete(participant).Error

	return err
}
