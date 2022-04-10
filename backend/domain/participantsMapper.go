package domain

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantDAO struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
}

func (ParticipantDAO) TableName() string {
	return "participants"
}

func ToParticipantDAO(participant *Participant) *ParticipantDAO {
	return &ParticipantDAO{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		CreatedAt:      participant.CreatedAt,
	}
}
