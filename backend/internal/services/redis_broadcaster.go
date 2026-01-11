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
)

type RedisBroadcaster interface {
	PublishNotification(ctx context.Context, notification ws.OutgoingNotification, serverID string) error
	PublishInvalidate(ctx context.Context, userID uuid.UUID) error
	Subscribe(ctx context.Context) error
	Close() error
}

type redisBroadcaster struct {
	redisClient         *redis.Client
	subscriptionChannel string
	chatChannel         string

	notificationService NotificationService
	chatPubsub          *redis.PubSub
	subPubsub           *redis.PubSub
}

func NewRedisBroadcaster(redisClient *redis.Client, notificationService NotificationService) RedisBroadcaster {
	return &redisBroadcaster{
		redisClient:         redisClient,
		subscriptionChannel: pubsub.SubscriptionChannel,
		chatChannel:         pubsub.ChatChannel,
		notificationService: notificationService,
	}
}

func (b *redisBroadcaster) PublishNotification(ctx context.Context, notification ws.OutgoingNotification, serverID string) error {
	messageID := uuid.New().String()

	bMessage := BroadcastMessage{
		Payload:        notification,
		UserID:         notification.UserID,
		MessageID:      messageID,
		ServerID:       serverID,
		ConversationID: uuid.Nil,
	}

	data, err := json.Marshal(bMessage)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := b.redisClient.Publish(ctx, b.chatChannel, data).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (b *redisBroadcaster) PublishInvalidate(ctx context.Context, userID uuid.UUID) error {
	event := SubscriptionEvent{
		Action: "invalidate",
		UserID: userID,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := b.redisClient.Publish(ctx, b.subscriptionChannel, data).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (b *redisBroadcaster) Subscribe(ctx context.Context) error {
	chatPubsub := b.redisClient.Subscribe(ctx, b.chatChannel)
	b.chatPubsub = chatPubsub

	subPubsub := b.redisClient.Subscribe(ctx, b.subscriptionChannel)
	b.subPubsub = subPubsub

	chatCh := chatPubsub.Channel()
	subCh := subPubsub.Channel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Redis subscription goroutine recovered from panic: %v", r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-chatCh:
				if !ok {
					return
				}

				var broadcastMsg BroadcastMessage
				if err := json.Unmarshal([]byte(msg.Payload), &broadcastMsg); err != nil {
					log.Printf("Error unmarshaling broadcast message: %v, payload: %s", err, msg.Payload)
					continue
				}

				if err := b.notificationService.Broadcast(ctx, broadcastMsg.ConversationID, broadcastMsg.Payload); err != nil {
					log.Printf("Error broadcasting message: %v", err)
				}
			case msg, ok := <-subCh:
				if !ok {
					return
				}

				var event SubscriptionEvent
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					log.Printf("Error unmarshaling subscription event: %v, payload: %s", err, msg.Payload)
					continue
				}

				if event.Action == "invalidate" {
					if err := b.notificationService.InvalidateMembership(ctx, event.UserID); err != nil {
						log.Printf("Error invalidating membership: %v", err)
					}
				}
			}
		}
	}()

	return nil
}

func (b *redisBroadcaster) Close() error {
	if b.chatPubsub != nil {
		if err := b.chatPubsub.Close(); err != nil {
			log.Printf("Error closing chat pubsub: %v", err)
		}
	}
	if b.subPubsub != nil {
		if err := b.subPubsub.Close(); err != nil {
			log.Printf("Error closing sub pubsub: %v", err)
		}
	}
	return nil
}
