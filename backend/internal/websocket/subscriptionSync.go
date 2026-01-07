package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type SubscriptionEvent struct {
	Action    string    `json:"action"`
	ClientID  uuid.UUID `json:"client_id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
}

type SubscriptionSync interface {
	PublishSubscribe(clientID uuid.UUID, channelID uuid.UUID, userID uuid.UUID) error
	PublishUnsubscribe(clientID uuid.UUID, channelID uuid.UUID, userID uuid.UUID) error
	Run(ctx context.Context, handleEvent func(SubscriptionEvent))
}

type subscriptionSync struct {
	redisClient *redis.Client
	channel     string
}

func NewSubscriptionSync(redisClient *redis.Client, channel string) *subscriptionSync {
	return &subscriptionSync{
		redisClient: redisClient,
		channel:     channel,
	}
}

func (ss *subscriptionSync) PublishSubscribe(clientID uuid.UUID, channelID uuid.UUID, userID uuid.UUID) error {
	event := SubscriptionEvent{
		Action:    "subscribe",
		ClientID:  clientID,
		ChannelID: channelID,
		UserID:    userID,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := ss.redisClient.Publish(context.Background(), ss.channel, data).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (ss *subscriptionSync) PublishUnsubscribe(clientID uuid.UUID, channelID uuid.UUID, userID uuid.UUID) error {
	event := SubscriptionEvent{
		Action:    "unsubscribe",
		ClientID:  clientID,
		ChannelID: channelID,
		UserID:    userID,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := ss.redisClient.Publish(context.Background(), ss.channel, data).Err(); err != nil {
		return fmt.Errorf("redis publish error: %w", err)
	}

	return nil
}

func (ss *subscriptionSync) Run(ctx context.Context, handleEvent func(SubscriptionEvent)) {
	redisPubsub := ss.redisClient.Subscribe(ctx, ss.channel)
	subscriptionChannel := redisPubsub.Channel()
	defer func() {
		_ = redisPubsub.Close()
	}()

	for {
		select {
		case message := <-subscriptionChannel:
			var event SubscriptionEvent

			if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
				continue
			}

			handleEvent(event)

		case <-ctx.Done():
			return
		}
	}
}
