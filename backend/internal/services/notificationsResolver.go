package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type NotificationResolverService interface {
	GetConversationRecipients(conversationId uuid.UUID) ([]uuid.UUID, error)
}

type notificationResolverService struct {
	participants domain.ParticipantRepository
}

func NewNotificationResolverService(
	participants domain.ParticipantRepository,
) *notificationResolverService {
	return &notificationResolverService{
		participants: participants,
	}
}

func (s *notificationResolverService) GetConversationRecipients(conversationID uuid.UUID) ([]uuid.UUID, error) {
	return s.participants.GetIDsByConversationID(conversationID)
}
