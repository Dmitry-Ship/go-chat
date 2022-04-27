package services

import (
	"context"
	"encoding/json"
	"log"

	"GitHub/go-chat/backend/internal/domain"
	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/infra/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type NotificationService interface {
	SubscribeToTopic(topic string, userId uuid.UUID) error
	BroadcastToTopic(topic string, notification ws.OutgoingNotification) error
	UnsubscribeFromTopic(topic string, userId uuid.UUID) error
	DeleteTopic(topic string) error
}

type NotificationClientRegister interface {
	RegisterClient(conn *websocket.Conn, wsHandlers ws.WSHandlers, userID uuid.UUID) error
}

type broadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserIDs []uuid.UUID             `json:"user_ids"`
}

type subscrition struct {
	userID uuid.UUID
	Topic  string
}

type notificationService struct {
	registerClient     chan *ws.Client
	unregisterClient   chan *ws.Client
	userClientsMap     map[uuid.UUID]map[uuid.UUID]*ws.Client
	redisClient        *redis.Client
	notificationTopics domain.NotificationTopicCommandRepository
}

var ctx = context.Background()

func NewNotificationService(redisClient *redis.Client, notificationTopics domain.NotificationTopicCommandRepository) *notificationService {
	return &notificationService{
		registerClient:     make(chan *ws.Client),
		unregisterClient:   make(chan *ws.Client),
		userClientsMap:     make(map[uuid.UUID]map[uuid.UUID]*ws.Client),
		redisClient:        redisClient,
		notificationTopics: notificationTopics,
	}
}

func (s *notificationService) Run() {
	redisPubsub := s.redisClient.Subscribe(ctx, pubsub.ChatChannel)
	ch := redisPubsub.Channel()
	defer redisPubsub.Close()

	for {
		select {
		case message := <-ch:
			if message.Payload == "ping" {
				s.redisClient.Publish(ctx, pubsub.ChatChannel, "pong")
				continue
			}

			var bMessage broadcastMessage

			err := json.Unmarshal([]byte(message.Payload), &bMessage)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, userID := range bMessage.UserIDs {
				clients, ok := s.userClientsMap[userID]

				if !ok {
					continue
				}

				for _, client := range clients {
					client.SendNotification(&bMessage.Payload)
				}
			}

		case client := <-s.registerClient:
			userClients, ok := s.userClientsMap[client.UserID]

			if !ok {
				userClients = make(map[uuid.UUID]*ws.Client)
				s.userClientsMap[client.UserID] = userClients
			}

			userClients[client.Id] = client

		case client := <-s.unregisterClient:
			if _, ok := s.userClientsMap[client.UserID]; ok {
				delete(s.userClientsMap[client.UserID], client.Id)
				close(client.SendChannel)
			}
		}
	}
}

func (s *notificationService) broadcastToUsers(userIDs []uuid.UUID, notification ws.OutgoingNotification) error {
	message := broadcastMessage{
		Payload: notification,
		UserIDs: userIDs,
	}

	json, err := json.Marshal(message)

	if err != nil {
		return err
	}

	s.redisClient.Publish(ctx, pubsub.ChatChannel, []byte(json))

	return nil
}

func (s *notificationService) RegisterClient(conn *websocket.Conn, wsHandlers ws.WSHandlers, userID uuid.UUID) error {
	client := ws.NewClient(conn, s.unregisterClient, wsHandlers, userID)

	go client.WritePump()
	go client.ReadPump()

	s.registerClient <- client

	return nil
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

	err = s.broadcastToUsers(userIds, notification)

	if err != nil {
		return err
	}

	return nil
}
