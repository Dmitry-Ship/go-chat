package services

import (
	"GitHub/go-chat/backend/pkg/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"github.com/google/uuid"
)

type NotificationsService interface {
	NotifyAboutMessage(conversationId uuid.UUID, messageID uuid.UUID, userId uuid.UUID) error
	NotifyAboutConversationDeletion(conversationId uuid.UUID)
	NotifyAboutConversationRenamed(conversationId uuid.UUID, newName string)
	SubscribeToTopic(topic string, userId uuid.UUID) error
	UnsubscribeFromTopic(topic string, userId uuid.UUID) error
	RegisterClient(client *ws.Client) error
}

type notificationsService struct {
	messages           readModel.MessageQueryRepository
	notificationTopics domain.NotificationTopicCommandRepository
	conencionsHub      ws.Hub
}

func NewNotificationsService(
	messages readModel.MessageQueryRepository,
	notificationTopics domain.NotificationTopicCommandRepository,
	conencionsHub ws.Hub,
) *notificationsService {
	return &notificationsService{
		messages:           messages,
		notificationTopics: notificationTopics,
		conencionsHub:      conencionsHub,
	}
}

func (s *notificationsService) RegisterClient(client *ws.Client) error {
	s.conencionsHub.RegisterClient(client)

	notificationTopics, err := s.notificationTopics.GetAllNotificationTopics(client.UserID)

	if err != nil {
		return err
	}

	go s.conencionsHub.SubscribeToTopic("user:"+client.UserID.String(), client.UserID)

	for _, topic := range notificationTopics {
		go s.conencionsHub.SubscribeToTopic(topic, client.UserID)
	}

	return nil
}

func (s *notificationsService) SubscribeToTopic(topic string, userId uuid.UUID) error {
	s.conencionsHub.SubscribeToTopic(topic, userId)

	notificationTopic := domain.NewNotificationTopic(topic, userId)

	err := s.notificationTopics.Store(notificationTopic)

	return err
}

func (s *notificationsService) UnsubscribeFromTopic(topic string, userId uuid.UUID) error {
	s.conencionsHub.UnsubscribeFromTopic(topic, userId)

	err := s.notificationTopics.DeleteByUserIDAndTopic(userId, topic)

	return err
}

func (s *notificationsService) notifyParticipants(conversationID uuid.UUID, notification ws.OutgoingNotification) {
	s.conencionsHub.BroadcastToTopic("conversation:"+conversationID.String(), notification)
}

func (s *notificationsService) NotifyAboutMessage(conversationId uuid.UUID, messageID uuid.UUID, userId uuid.UUID) error {
	messageDTO, err := s.messages.GetMessageByID(messageID, userId)

	if err != nil {
		return err
	}

	notification := ws.OutgoingNotification{
		Type:    "message",
		Payload: messageDTO,
	}

	s.notifyParticipants(conversationId, notification)

	return nil
}

func (s *notificationsService) NotifyAboutConversationDeletion(id uuid.UUID) {
	notification := ws.OutgoingNotification{
		Type: "conversation_deleted",
		Payload: struct {
			ConversationId uuid.UUID `json:"conversation_id"`
		}{
			ConversationId: id,
		},
	}

	s.notifyParticipants(id, notification)
}

func (s *notificationsService) NotifyAboutConversationRenamed(conversationId uuid.UUID, newName string) {
	notification := ws.OutgoingNotification{
		Type: "conversation_renamed",
		Payload: struct {
			ConversationId uuid.UUID `json:"conversation_id"`
			NewName        string    `json:"new_name"`
		}{
			ConversationId: conversationId,
			NewName:        newName,
		},
	}

	s.notifyParticipants(conversationId, notification)
}