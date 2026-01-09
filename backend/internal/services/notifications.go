package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"GitHub/go-chat/backend/internal/domain"
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

type SubscriptionEvent struct {
	Action string    `json:"action"`
	UserID uuid.UUID `json:"user_id"`
}

type NotificationService interface {
	Broadcast(ctx context.Context, channelID uuid.UUID, notification ws.OutgoingNotification) error
	RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID
	Run()
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
}

type notificationService struct {
	ctx                 context.Context
	activeClients       ws.ActiveClients
	redisClient         *redis.Client
	subscriptionChannel string
}

func NewNotificationService(
	ctx context.Context,
	redisClient *redis.Client,
	participants domain.ParticipantRepository,
) *notificationService {
	return &notificationService{
		ctx:                 ctx,
		activeClients:       ws.NewActiveClients(participants),
		redisClient:         redisClient,
		subscriptionChannel: pubsub.SubscriptionChannel,
	}
}

func NewNotificationServiceWithClients(
	ctx context.Context,
	redisClient *redis.Client,
	activeClients ws.ActiveClients,
) *notificationService {
	return &notificationService{
		ctx:                 ctx,
		activeClients:       activeClients,
		redisClient:         redisClient,
		subscriptionChannel: pubsub.SubscriptionChannel,
	}
}

func (s *notificationService) Broadcast(ctx context.Context, channelID uuid.UUID, notification ws.OutgoingNotification) error {
	if channelID != uuid.Nil {
		clients := s.activeClients.GetClientsByChannel(channelID)
		for _, client := range clients {
			client.SendNotification(notification)
		}
	}

	bMessage := BroadcastMessage{
		Payload: notification,
		UserID:  notification.UserID,
	}

	data, err := json.Marshal(bMessage)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := s.redisClient.Publish(ctx, pubsub.ChatChannel, data).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (s *notificationService) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)

	clientID := s.activeClients.AddClient(newClient)

	if err := s.activeClients.InvalidateMembership(ctx, userID); err != nil {
		log.Printf("Error invalidating membership: %v", err)
	}

	go newClient.WritePump()
	newClient.ReadPump()

	return clientID
}

func (s *notificationService) Run() {
	redisPubsub := s.redisClient.Subscribe(s.ctx, pubsub.ChatChannel)
	chatChannel := redisPubsub.Channel()
	defer func() {
		_ = redisPubsub.Close()
	}()

	subscriptionPubsub := s.redisClient.Subscribe(s.ctx, s.subscriptionChannel)
	subscriptionChannel := subscriptionPubsub.Channel()
	defer func() {
		_ = subscriptionPubsub.Close()
	}()

	go func() {
		for {
			select {
			case message := <-subscriptionChannel:
				var event SubscriptionEvent

				if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
					continue
				}

				switch event.Action {
				case "invalidate":
					if err := s.activeClients.InvalidateMembership(s.ctx, event.UserID); err != nil {
						log.Printf("Error invalidating membership: %v", err)
					}
				}

			case <-s.ctx.Done():
				return
			}
		}
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

			clients := s.activeClients.GetClientsByUser(bMessage.UserID)
			for _, client := range clients {
				client.SendNotification(bMessage.Payload)
			}

		case <-s.ctx.Done():
			return
		}
	}
}

func (s *notificationService) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	if err := s.activeClients.InvalidateMembership(ctx, userID); err != nil {
		return err
	}

	if err := s.publishInvalidate(userID); err != nil {
		log.Printf("Error publishing invalidate: %v", err)
	}

	return nil
}

func (s *notificationService) publishInvalidate(userID uuid.UUID) error {
	event := SubscriptionEvent{
		Action: "invalidate",
		UserID: userID,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := s.redisClient.Publish(s.ctx, s.subscriptionChannel, data).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}
