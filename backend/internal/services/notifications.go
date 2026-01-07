package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserID  uuid.UUID               `json:"user_id"`
}

type NotificationService interface {
	Send(message ws.OutgoingNotification) error
	RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte))
	Run()
}

type notificationService struct {
	ctx           context.Context
	activeClients ws.ActiveClients
	redisClient   *redis.Client
}

func NewNotificationService(
	ctx context.Context,
	redisClient *redis.Client,
) *notificationService {
	return &notificationService{
		ctx:           ctx,
		activeClients: ws.NewActiveClients(),
		redisClient:   redisClient,
	}
}

func (s *notificationService) Send(message ws.OutgoingNotification) error {
	bMessage := BroadcastMessage{
		Payload: message,
		UserID:  message.UserID,
	}

	json, err := json.Marshal(bMessage)

	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (s *notificationService) RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)

	s.activeClients.AddClient(newClient)

	go newClient.WritePump()
	newClient.ReadPump()
}

func (s *notificationService) Run() {
	redisPubsub := s.redisClient.Subscribe(s.ctx, pubsub.ChatChannel)
	chatChannel := redisPubsub.Channel()
	defer func() {
		_ = redisPubsub.Close()
	}()

	for {
		select {
		case message := <-chatChannel:
			if message.Payload == "ping" {
				s.redisClient.Publish(s.ctx, pubsub.ChatChannel, "pong")
				continue
			}

			var bMessage BroadcastMessage

			if err := json.Unmarshal([]byte(message.Payload), &bMessage); err != nil {
				log.Println(err)
				continue
			}

			s.activeClients.SendToUserClients(bMessage.UserID, bMessage.Payload)

		case <-s.ctx.Done():
			return
		}
	}
}
