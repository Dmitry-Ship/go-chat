package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type PublicConversationParticipationService interface {
	Join(conversationId uuid.UUID, userId uuid.UUID) error
	Leave(conversationId uuid.UUID, userId uuid.UUID) error
}

type publicConversationParticipationService struct {
	participants domain.ParticipantRepository
}

func NewPublicConversationParticipationService(
	participants domain.ParticipantRepository,
) *publicConversationParticipationService {
	return &publicConversationParticipationService{
		participants: participants,
	}
}

func (s *publicConversationParticipationService) Join(conversationID uuid.UUID, userId uuid.UUID) error {
	return s.participants.Store(domain.NewJoinedParticipant(conversationID, userId))
}

func (s *publicConversationParticipationService) Leave(conversationID uuid.UUID, userId uuid.UUID) error {
	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	err = participant.LeavePublicConversation(conversationID)

	if err != nil {
		return err
	}

	return s.participants.Update(participant)
}
