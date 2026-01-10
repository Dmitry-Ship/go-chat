package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

type BatchBroadcastMessage struct {
	Payload ws.BatchOutgoingNotification `json:"notification"`
	UserID  uuid.UUID                    `json:"user_id"`
}

type SubscriptionEvent struct {
	Action string    `json:"action"`
	UserID uuid.UUID `json:"user_id"`
}

type batchKey struct {
	channelID uuid.UUID
	userID    uuid.UUID
}

type batchNotification struct {
	channelID    uuid.UUID
	notification ws.OutgoingNotification
}

type notificationService struct {
	ctx                 context.Context
	activeClients       ws.ActiveClients
	redisClient         *redis.Client
	subscriptionChannel string
	batchChannel        chan batchNotification
	batchTimeout        time.Duration
	maxBatchSize        int
}

type NotificationService interface {
	Broadcast(ctx context.Context, channelID uuid.UUID, notification ws.OutgoingNotification) error
	RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID, handleNotification func(userID uuid.UUID, message []byte)) uuid.UUID
	Run()
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
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
		batchChannel:        make(chan batchNotification, 1000),
		batchTimeout:        10 * time.Millisecond,
		maxBatchSize:        10,
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
		batchChannel:        make(chan batchNotification, 1000),
		batchTimeout:        10 * time.Millisecond,
		maxBatchSize:        10,
	}
}

func (s *notificationService) Broadcast(ctx context.Context, channelID uuid.UUID, notification ws.OutgoingNotification) error {
	s.batchChannel <- batchNotification{
		channelID:    channelID,
		notification: notification,
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
	go s.runBatcher()

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

func (s *notificationService) runBatcher() {
	batches := make(map[batchKey][]ws.OutgoingNotification)
	timer := time.NewTimer(s.batchTimeout)

	for {
		select {
		case bn := <-s.batchChannel:
			key := batchKey{channelID: bn.channelID, userID: bn.notification.UserID}
			batches[key] = append(batches[key], bn.notification)
			if len(batches[key]) >= s.maxBatchSize {
				s.flushBatch(key, batches[key])
				delete(batches, key)
			} else {
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(s.batchTimeout)
			}

		case <-timer.C:
			for key, notifications := range batches {
				s.flushBatch(key, notifications)
				delete(batches, key)
			}
			timer.Reset(s.batchTimeout)

		case <-s.ctx.Done():
			for key, notifications := range batches {
				s.flushBatch(key, notifications)
			}
			return
		}
	}
}

func (s *notificationService) flushBatch(key batchKey, notifications []ws.OutgoingNotification) {
	if len(notifications) == 0 {
		return
	}

	if len(notifications) == 1 {
		notification := notifications[0]
		if key.channelID != uuid.Nil {
			clients := s.activeClients.GetClientsByChannel(key.channelID)
			for _, client := range clients {
				client.SendNotification(notification)
			}
		}
		clients := s.activeClients.GetClientsByUser(notification.UserID)
		for _, client := range clients {
			client.SendNotification(notification)
		}

		bMessage := BroadcastMessage{
			Payload: notification,
			UserID:  notification.UserID,
		}
		data, err := json.Marshal(bMessage)
		if err != nil {
			log.Printf("json marshal error: %v", err)
			return
		}
		if err := s.redisClient.Publish(s.ctx, pubsub.ChatChannel, data).Err(); err != nil {
			log.Printf("redis publish error: %v", err)
		}
		return
	}

	batch := ws.BatchOutgoingNotification{
		UserID: notifications[0].UserID,
		Events: make([]ws.NotificationEvent, len(notifications)),
	}
	for i, n := range notifications {
		batch.Events[i] = ws.NotificationEvent{
			Type:    n.Type,
			Payload: n.Payload,
		}
	}

	if key.channelID != uuid.Nil {
		clients := s.activeClients.GetClientsByChannel(key.channelID)
		for _, client := range clients {
			client.SendNotification(ws.OutgoingNotification{
				Type:    "batch",
				UserID:  batch.UserID,
				Payload: batch,
			})
		}
	}

	clients := s.activeClients.GetClientsByUser(batch.UserID)
	for _, client := range clients {
		client.SendNotification(ws.OutgoingNotification{
			Type:    "batch",
			UserID:  batch.UserID,
			Payload: batch,
		})
	}

	for _, n := range notifications {
		bMessage := BroadcastMessage{
			Payload: n,
			UserID:  n.UserID,
		}
		data, err := json.Marshal(bMessage)
		if err != nil {
			log.Printf("json marshal error: %v", err)
			continue
		}
		if err := s.redisClient.Publish(s.ctx, pubsub.ChatChannel, data).Err(); err != nil {
			log.Printf("redis publish error: %v", err)
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
