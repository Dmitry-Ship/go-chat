package services

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/infra/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type NotificationsService interface {
	SubscribeToTopic(topic string, userId uuid.UUID) error
	BroadcastToTopic(topic string, notification ws.OutgoingNotification)
	UnsubscribeFromTopic(topic string, userId uuid.UUID) error
	RegisterClient(conn *websocket.Conn, wsHandlers ws.WSHandlers, userID uuid.UUID) error
	DeleteTopic(topic string) error
}

type notificationsService struct {
	notificationTopics domain.NotificationTopicCommandRepository
	conencionsHub      ws.Hub
}

func NewNotificationsService(
	notificationTopics domain.NotificationTopicCommandRepository,
	conencionsHub ws.Hub,
) *notificationsService {
	return &notificationsService{
		notificationTopics: notificationTopics,
		conencionsHub:      conencionsHub,
	}
}

func (s *notificationsService) RegisterClient(conn *websocket.Conn, wsHandlers ws.WSHandlers, userID uuid.UUID) error {
	client := ws.NewClient(conn, s.conencionsHub, wsHandlers, userID)

	go client.WritePump()
	go client.ReadPump()

	s.conencionsHub.RegisterClient(client)

	notificationTopics, err := s.notificationTopics.GetAllNotificationTopics(client.UserID)

	if err != nil {
		return err
	}

	for _, topic := range notificationTopics {
		go s.conencionsHub.SubscribeToTopic(topic, client.UserID)
	}

	go s.conencionsHub.SubscribeToTopic("user:"+client.UserID.String(), client.UserID)

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

func (s *notificationsService) DeleteTopic(topic string) error {
	s.conencionsHub.DeleteTopic(topic)

	return s.notificationTopics.DeleteByTopic(topic)
}

func (s *notificationsService) BroadcastToTopic(topic string, notification ws.OutgoingNotification) {
	s.conencionsHub.BroadcastToTopic(topic, notification)
}
