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
	SubscribeToChat(ctx context.Context) <-chan BroadcastMessage
	SubscribeToInvalidation(ctx context.Context) <-chan SubscriptionEvent
	Close()
}

type redisBroadcaster struct {
	redisClient         *redis.Client
	subscriptionChannel string
	chatChannel         string
	chatMessages        chan BroadcastMessage
	invalidationEvents  chan SubscriptionEvent
	chatPubsub          *redis.PubSub
	subPubsub           *redis.PubSub
}

func NewRedisBroadcaster(redisClient *redis.Client) RedisBroadcaster {
	return &redisBroadcaster{
		redisClient:         redisClient,
		subscriptionChannel: pubsub.SubscriptionChannel,
		chatChannel:         pubsub.ChatChannel,
		chatMessages:        make(chan BroadcastMessage, 1000),
		invalidationEvents:  make(chan SubscriptionEvent, 1000),
	}
}

func (b *redisBroadcaster) PublishNotification(ctx context.Context, notification ws.OutgoingNotification, serverID string) error {
	messageID := uuid.New().String()

	bMessage := BroadcastMessage{
		Payload:   notification,
		UserID:    notification.UserID,
		MessageID: messageID,
		ServerID:  serverID,
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

func (b *redisBroadcaster) SubscribeToChat(ctx context.Context) <-chan BroadcastMessage {
	b.chatPubsub = b.redisClient.Subscribe(ctx, b.chatChannel)
	chatChannel := b.chatPubsub.Channel()

	go b.readChatMessages(ctx, chatChannel)
	return b.chatMessages
}

func (b *redisBroadcaster) SubscribeToInvalidation(ctx context.Context) <-chan SubscriptionEvent {
	b.subPubsub = b.redisClient.Subscribe(ctx, b.subscriptionChannel)
	subscriptionChannel := b.subPubsub.Channel()

	go b.readInvalidationEvents(ctx, subscriptionChannel)
	return b.invalidationEvents
}

func (b *redisBroadcaster) readChatMessages(ctx context.Context, chatChannel <-chan *redis.Message) {
	for {
		select {
		case message := <-chatChannel:
			if message == nil {
				return
			}

			var bMessage BroadcastMessage
			if err := json.Unmarshal([]byte(message.Payload), &bMessage); err != nil {
				log.Printf("Error unmarshaling broadcast message: %v, payload: %s", err, message.Payload)
				continue
			}

			select {
			case b.chatMessages <- bMessage:
			case <-ctx.Done():
				return
			default:
				log.Printf("Chat messages channel full, dropping message: %s", bMessage.MessageID)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (b *redisBroadcaster) readInvalidationEvents(ctx context.Context, subscriptionChannel <-chan *redis.Message) {
	for {
		select {
		case message := <-subscriptionChannel:
			if message == nil {
				return
			}

			var event SubscriptionEvent
			if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
				log.Printf("Error unmarshaling subscription event: %v, payload: %s", err, message.Payload)
				continue
			}

			select {
			case b.invalidationEvents <- event:
			case <-ctx.Done():
				return
			default:
				log.Printf("Invalidation events channel full, dropping event for user: %s", event.UserID)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (b *redisBroadcaster) Close() {
	if b.chatPubsub != nil {
		_ = b.chatPubsub.Close()
	}
	if b.subPubsub != nil {
		_ = b.subPubsub.Close()
	}
	close(b.chatMessages)
	close(b.invalidationEvents)
}
