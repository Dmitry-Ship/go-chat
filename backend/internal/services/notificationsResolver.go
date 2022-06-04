package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type NotificationResolverService interface {
	GetConversationRecipients(conversationId uuid.UUID) ([]uuid.UUID, error)
	GetReceiversFromEvent(event domain.DomainEvent) ([]uuid.UUID, error)
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

func (s *notificationResolverService) GetReceiversFromEvent(event domain.DomainEvent) ([]uuid.UUID, error) {
	var receiversIDs []uuid.UUID
	var err error

	// find notification receivers
	switch e := event.(type) {
	case
		*domain.GroupConversationRenamed,
		*domain.GroupConversationLeft,
		*domain.GroupConversationJoined,
		*domain.GroupConversationInvited,
		*domain.MessageSent,
		*domain.GroupConversationDeleted:
		if e, ok := e.(domain.ConversationEvent); ok {
			receiversIDs, err = s.GetConversationRecipients(e.GetConversationID())

			if err != nil {
				return receiversIDs, err
			}
		}
	}

	return receiversIDs, nil
}
