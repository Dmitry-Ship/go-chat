package services

import (
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type NotificationBuilderService interface {
	GetConversationDeletedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (*ws.OutgoingNotification, error)
	GetConversationUpdatedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (*ws.OutgoingNotification, error)
	GetMessageSentBuilder(messageID uuid.UUID) func(userID uuid.UUID) (*ws.OutgoingNotification, error)
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

func (s *notificationBuilderService) GetConversationUpdatedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		conversation, err := s.queries.GetConversation(conversationID, userID)

		if err != nil {
			return nil, err
		}

		return &ws.OutgoingNotification{
			Type:    "conversation_updated",
			Payload: conversation,
		}, nil
	}
}

func (s *notificationBuilderService) GetConversationDeletedBuilder(conversationID uuid.UUID) func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		return &ws.OutgoingNotification{
			Type: "conversation_deleted",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: conversationID,
			},
		}, nil
	}
}

func (s *notificationBuilderService) GetMessageSentBuilder(messageID uuid.UUID) func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
	return func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		messageDTO, err := s.queries.GetNotificationMessage(messageID, userID)

		if err != nil {
			return nil, err
		}

		return &ws.OutgoingNotification{
			Type:    "message",
			Payload: messageDTO,
		}, nil
	}
}
