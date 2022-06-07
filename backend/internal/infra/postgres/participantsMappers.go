package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

func toParticipantPersistence(participant domain.Participant) *Participant {
	return &Participant{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		IsActive:       participant.IsActive,
	}
}

func toParticipantDomain(participant Participant) *domain.Participant {
	return &domain.Participant{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		IsActive:       participant.IsActive,
	}
}
