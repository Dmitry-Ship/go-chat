package hub

import (
	"context"
	"encoding/json"
	"log"

	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
)

type BroadcastMessage struct {
	Payload ws.OutgoingNotification `json:"notification"`
	UserIDs []uuid.UUID             `json:"user_ids"`
}

type broadcaster struct {
	ctx           context.Context
	activeClients ws.ActiveClients
	redisClient   *redis.Client
}

func NewBroadcaster(ctx context.Context, redisClient *redis.Client, activeClients ws.ActiveClients) *broadcaster {
	return &broadcaster{
		ctx:           ctx,
		activeClients: activeClients,
		redisClient:   redisClient,
	}
}

func (s *broadcaster) Run() {
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

			var bMessage BroadcastMessage

			err := json.Unmarshal([]byte(message.Payload), &bMessage)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, userID := range bMessage.UserIDs {
				s.activeClients.SendToUserClients(userID, bMessage.Payload)
			}
		}
	}
}
