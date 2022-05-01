package services

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/infra/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type NotificationService interface {
	SubscribeToTopic(topic string, userId uuid.UUID) error
	BroadcastToTopic(topic string, notification ws.OutgoingNotification) error
	UnsubscribeFromTopic(topic string, userId uuid.UUID) error
	DeleteTopic(topic string) error
	RegisterClient(conn *websocket.Conn, userID uuid.UUID) error
}

type broadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserIDs []uuid.UUID             `json:"user_ids"`
}

type notificationService struct {
	connectionsPool    ws.ConnectionsPool
	notificationTopics domain.NotificationTopicCommandRepository
}

func NewNotificationService(connectionsPool ws.ConnectionsPool, notificationTopics domain.NotificationTopicCommandRepository) *notificationService {
	return &notificationService{
		connectionsPool:    connectionsPool,
		notificationTopics: notificationTopics,
	}
}

func (s *notificationService) RegisterClient(conn *websocket.Conn, userID uuid.UUID) error {
	return s.connectionsPool.RegisterClient(conn, userID)
}

func (s *notificationService) SubscribeToTopic(topic string, userId uuid.UUID) error {
	notificationTopic := domain.NewNotificationTopic(topic, userId)

	return s.notificationTopics.Store(notificationTopic)
}

func (s *notificationService) UnsubscribeFromTopic(topic string, userId uuid.UUID) error {
	return s.notificationTopics.DeleteByUserIDAndTopic(userId, topic)
}

func (s *notificationService) DeleteTopic(topic string) error {
	return s.notificationTopics.DeleteAllByTopic(topic)
}

func (s *notificationService) BroadcastToTopic(topic string, notification ws.OutgoingNotification) error {
	userIds, err := s.notificationTopics.GetUserIDsByTopic(topic)

	if err != nil {
		return err
	}

	return s.connectionsPool.BroadcastToUsers(userIds, notification)
}
