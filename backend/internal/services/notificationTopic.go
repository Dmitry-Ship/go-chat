package services

import (
	"GitHub/go-chat/backend/internal/domain"
	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type BroadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserID  uuid.UUID               `json:"user_id"`
}

type NotificationTopicService interface {
	SubscribeToTopic(topic string, userID uuid.UUID) error
	GetReceivers(topic string) ([]uuid.UUID, error)
	SendToUser(userID uuid.UUID, notification ws.OutgoingNotification) error
	UnsubscribeFromTopic(topic string, userID uuid.UUID) error
	DeleteTopic(topic string) error
}

type notificationTopicService struct {
	ctx                context.Context
	notificationTopics domain.NotificationTopicRepository
	redisClient        *redis.Client
}

func NewNotificationTopicService(
	ctx context.Context,
	notificationTopics domain.NotificationTopicRepository,
	redisClient *redis.Client,
) *notificationTopicService {
	return &notificationTopicService{
		ctx:                ctx,
		notificationTopics: notificationTopics,
		redisClient:        redisClient,
	}
}

func (s *notificationTopicService) SubscribeToTopic(topic string, userID uuid.UUID) error {
	notificationTopicID := uuid.New()

	notificationTopic, err := domain.NewNotificationTopic(notificationTopicID, topic, userID)

	if err != nil {
		return err
	}

	return s.notificationTopics.Store(notificationTopic)
}

func (s *notificationTopicService) UnsubscribeFromTopic(topic string, userID uuid.UUID) error {
	return s.notificationTopics.DeleteByUserIDAndTopic(userID, topic)
}

func (s *notificationTopicService) DeleteTopic(topic string) error {
	return s.notificationTopics.DeleteAllByTopic(topic)
}

func (s *notificationTopicService) GetReceivers(topic string) ([]uuid.UUID, error) {
	return s.notificationTopics.GetUserIDsByTopic(topic)
}

func (s *notificationTopicService) SendToUser(userID uuid.UUID, notification ws.OutgoingNotification) error {
	message := BroadcastMessage{
		Payload: notification,
		UserID:  userID,
	}

	json, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err()
}
