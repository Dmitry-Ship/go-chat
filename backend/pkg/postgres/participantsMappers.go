package postgres

import (
	"GitHub/go-chat/backend/domain"
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	ConversationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time
}

func toParticipantPersistence(participant *domain.ParticipantAggregate) *Participant {
	return &Participant{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		CreatedAt:      participant.CreatedAt,
	}
}
