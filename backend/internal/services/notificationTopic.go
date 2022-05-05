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
	UserIDs []uuid.UUID             `json:"user_ids"`
}

type NotificationTopicService interface {
	SubscribeToTopic(topic string, userId uuid.UUID) error
	BroadcastToTopic(topic string, notification ws.OutgoingNotification) error
	UnsubscribeFromTopic(topic string, userId uuid.UUID) error
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

func (s *notificationTopicService) SubscribeToTopic(topic string, userId uuid.UUID) error {
	notificationTopic := domain.NewNotificationTopic(topic, userId)

	return s.notificationTopics.Store(notificationTopic)
}

func (s *notificationTopicService) UnsubscribeFromTopic(topic string, userId uuid.UUID) error {
	return s.notificationTopics.DeleteByUserIDAndTopic(userId, topic)
}

func (s *notificationTopicService) DeleteTopic(topic string) error {
	return s.notificationTopics.DeleteAllByTopic(topic)
}

func (s *notificationTopicService) BroadcastToTopic(topic string, notification ws.OutgoingNotification) error {
	userIds, err := s.notificationTopics.GetUserIDsByTopic(topic)

	if err != nil {
		return err
	}

	message := BroadcastMessage{
		Payload: notification,
		UserIDs: userIds,
	}

	json, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err()
}
