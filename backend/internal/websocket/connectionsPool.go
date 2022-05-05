package ws

import (
	"context"
	"encoding/json"
	"log"

	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	pubsub "GitHub/go-chat/backend/internal/infra/redis"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type ClientRegister interface {
	RegisterClient(conn *websocket.Conn, userID uuid.UUID)
}

type NotificationTopicService interface {
	SubscribeToTopic(topic string, userId uuid.UUID) error
	BroadcastToTopic(topic string, notification OutgoingNotification) error
	UnsubscribeFromTopic(topic string, userId uuid.UUID) error
	DeleteTopic(topic string) error
}

type broadcastMessage struct {
	Payload OutgoingNotification `json:"notification"`
	UserIDs []uuid.UUID          `json:"user_ids"`
}

type hub struct {
	ctx                context.Context
	commands           *app.Commands
	notificationTopics domain.NotificationTopicRepository
	activeClients      *activeClients
	redisClient        *redis.Client
}

func NewHub(ctx context.Context, redisClient *redis.Client, commands *app.Commands, notificationTopics domain.NotificationTopicRepository) *hub {
	return &hub{
		ctx:                ctx,
		commands:           commands,
		notificationTopics: notificationTopics,
		activeClients:      newActiveClients(),
		redisClient:        redisClient,
	}
}

func (s *hub) Run() {
	redisPubsub := s.redisClient.Subscribe(s.ctx, pubsub.ChatChannel)
	ch := redisPubsub.Channel()
	defer redisPubsub.Close()

	for {
		select {
		case message := <-ch:
			if message.Payload == "ping" {
				s.redisClient.Publish(s.ctx, pubsub.ChatChannel, "pong")
				continue
			}

			var bMessage broadcastMessage

			err := json.Unmarshal([]byte(message.Payload), &bMessage)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, userID := range bMessage.UserIDs {
				s.activeClients.sendToUserClients(userID, bMessage.Payload)
			}
		}
	}
}

func (s *hub) RegisterClient(conn *websocket.Conn, userID uuid.UUID) {
	newClient := NewClient(conn, s.unregisterClient, s.handleNotification, userID)
	newClient.Listen()

	s.activeClients.addClient(newClient)
}

func (s *hub) unregisterClient(client *client) {
	s.activeClients.removeClient(client)
}

func (s *hub) handleNotification(notification *incomingNotification) {
	switch notification.Type {
	case "message":
		s.handleReceiveWSChatMessage(notification.Data, notification.UserID)
	default:
		log.Println("Unknown notification type:", notification.Type)
	}
}

func (s *hub) handleReceiveWSChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.commands.MessagingService.SendTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}

func (s *hub) broadcastToUsers(userIDs []uuid.UUID, notification OutgoingNotification) error {
	message := broadcastMessage{
		Payload: notification,
		UserIDs: userIDs,
	}

	json, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err()
}

func (s *hub) SubscribeToTopic(topic string, userId uuid.UUID) error {
	notificationTopic := domain.NewNotificationTopic(topic, userId)

	return s.notificationTopics.Store(notificationTopic)
}

func (s *hub) UnsubscribeFromTopic(topic string, userId uuid.UUID) error {
	return s.notificationTopics.DeleteByUserIDAndTopic(userId, topic)
}

func (s *hub) DeleteTopic(topic string) error {
	return s.notificationTopics.DeleteAllByTopic(topic)
}

func (s *hub) BroadcastToTopic(topic string, notification OutgoingNotification) error {
	userIds, err := s.notificationTopics.GetUserIDsByTopic(topic)

	if err != nil {
		return err
	}

	return s.broadcastToUsers(userIds, notification)
}
