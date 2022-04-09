package domain

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantPersistence struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
}

func (ParticipantPersistence) TableName() string {
	return "participants"
}

func ToParticipantPersistence(participant *Participant) *ParticipantPersistence {
	return &ParticipantPersistence{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		CreatedAt:      participant.CreatedAt,
	}
}
