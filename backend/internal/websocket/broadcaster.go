package ws

import (
	"context"
	"encoding/json"
	"log"

	pubsub "GitHub/go-chat/backend/internal/infra/redis"

	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
)

type BroadcastMessage struct {
	Payload OutgoingNotification `json:"notification"`
	UserID  uuid.UUID            `json:"user_id"`
}

type broadcaster struct {
	ctx           context.Context
	activeClients ActiveClients
	redisClient   *redis.Client
}

func NewBroadcaster(ctx context.Context, redisClient *redis.Client, activeClients ActiveClients) *broadcaster {
	return &broadcaster{
		ctx:           ctx,
		activeClients: activeClients,
		redisClient:   redisClient,
	}
}

func (s *broadcaster) Run() {
	redisPubsub := s.redisClient.Subscribe(s.ctx, pubsub.ChatChannel)
	chatChannel := redisPubsub.Channel()
	defer redisPubsub.Close()

	for {
		select {
		case message := <-chatChannel:
			if message.Payload == "ping" {
				s.redisClient.Publish(s.ctx, pubsub.ChatChannel, "pong")
				continue
			}

			var bMessage BroadcastMessage

			err := json.Unmarshal([]byte(message.Payload), &bMessage)

			if err != nil {
				log.Println(err)
				continue
			}

			s.activeClients.SendToUserClients(bMessage.UserID, bMessage.Payload)

		case <-s.ctx.Done():
			return
		}
	}
}
