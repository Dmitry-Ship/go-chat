package services

import (
	pubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type ClientsService interface {
	RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(notification *ws.IncomingNotification))
	Run()
}

type clientsService struct {
	activeClients ws.ActiveClients
	ctx           context.Context
	redisClient   *redis.Client
}

func NewClientsService(ctx context.Context, redisClient *redis.Client, activeClients ws.ActiveClients) *clientsService {
	return &clientsService{
		activeClients: activeClients,
		ctx:           ctx,
		redisClient:   redisClient,
	}
}

func (s *clientsService) RegisterClient(conn *websocket.Conn, userID uuid.UUID, handleNotification func(notification *ws.IncomingNotification)) {
	newClient := ws.NewClient(conn, s.activeClients.RemoveClient, handleNotification, userID)
	newClient.Listen()

	s.activeClients.AddClient(newClient)
}

func (s *clientsService) Run() {
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
