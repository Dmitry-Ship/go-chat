package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

var participantTypesMap = map[uint8]string{
	0: domain.ParticipantTypeJoined,
	1: domain.ParticipantTypeOwner,
	2: domain.ParticipantTypePrivate,
}

func toParticipantTypePersistence(participantType string) uint8 {
	for k, v := range participantTypesMap {
		if v == participantType {
			return k
		}
	}

	return 0
}

func toParticipantPersistence(participant *domain.Participant) *Participant {
	return &Participant{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		CreatedAt:      participant.CreatedAt,
		IsActive:       participant.IsActive,
		Type:           toParticipantTypePersistence(participant.Type),
	}
}

func toParticipantDomain(participant *Participant) *domain.Participant {
	return &domain.Participant{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		CreatedAt:      participant.CreatedAt,
		IsActive:       participant.IsActive,
		Type:           participantTypesMap[participant.Type],
	}
}
