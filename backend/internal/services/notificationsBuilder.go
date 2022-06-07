package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type NotificationBuilderService interface {
	GetConversationDeletedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error)
	GetConversationUpdatedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error)
	GetMessageSentBuilder(messageID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error)
	BuildMessageFromEvent(receiverID uuid.UUID, event domain.DomainEvent) (ws.OutgoingNotification, error)
}

type notificationBuilderService struct {
	queries readModel.QueriesRepository
}

func NewNotificationBuilderService(
	queries readModel.QueriesRepository,
) *notificationBuilderService {
	return &notificationBuilderService{
		queries: queries,
	}
}

func (s *notificationBuilderService) BuildMessageFromEvent(receiverID uuid.UUID, event domain.DomainEvent) (ws.OutgoingNotification, error) {
	var buildMessage func(userID uuid.UUID) (ws.OutgoingNotification, error)
	switch e := event.(type) {
	case domain.GroupConversationRenamed, domain.GroupConversationLeft, domain.GroupConversationJoined, domain.GroupConversationInvited:
		if e, ok := e.(domain.ConversationEvent); ok {
			buildMessage = s.GetConversationUpdatedBuilder(e.GetConversationID())
		}
	case domain.GroupConversationDeleted:
		buildMessage = s.GetConversationDeletedBuilder(e.GetConversationID())
	case domain.MessageSent:
		buildMessage = s.GetMessageSentBuilder(e.MessageID)
	}

	return buildMessage(receiverID)
}

func (s *notificationBuilderService) GetConversationUpdatedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (ws.OutgoingNotification, error) {
		conversation, err := s.queries.GetConversation(conversationID, userID)

		if err != nil {
			return ws.OutgoingNotification{}, err
		}

		return ws.OutgoingNotification{
			Type:    "conversation_updated",
			UserID:  userID,
			Payload: conversation,
		}, nil
	}
}

func (s *notificationBuilderService) GetConversationDeletedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (ws.OutgoingNotification, error) {
		return ws.OutgoingNotification{
			Type:   "conversation_deleted",
			UserID: userID,
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: conversationID,
			},
		}, nil
	}
}

func (s *notificationBuilderService) GetMessageSentBuilder(messageID uuid.UUID) func(userID uuid.UUID) (ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (ws.OutgoingNotification, error) {
		messageDTO, err := s.queries.GetNotificationMessage(messageID, userID)

		if err != nil {
			return ws.OutgoingNotification{}, err
		}

		return ws.OutgoingNotification{
			Type:    "message",
			UserID:  userID,
			Payload: messageDTO,
		}, nil
	}
}
