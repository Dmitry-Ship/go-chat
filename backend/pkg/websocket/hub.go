package ws

import (
	"context"
	"encoding/json"
	"log"

	pubsub "GitHub/go-chat/backend/pkg/redis"

	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
)

type Hub interface {
	BroadcastToClients(notification OutgoingNotification, recipientIDs []uuid.UUID)
	UnregisterClient(client *Client)
	RegisterClient(client *Client)
}

type broadcastMessage struct {
	Notification OutgoingNotification `json:"notification"`
	RecipientID  uuid.UUID            `json:"recipientID"`
}

type hub struct {
	register       chan *Client
	unregister     chan *Client
	userClientsMap map[uuid.UUID]map[uuid.UUID]*Client
	redisClient    *redis.Client
}

var ctx = context.Background()

func NewHub(redisClient *redis.Client) *hub {
	return &hub{
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		userClientsMap: make(map[uuid.UUID]map[uuid.UUID]*Client),
		redisClient:    redisClient,
	}
}

func (s *hub) Run() {
	redisPubsub := s.redisClient.Subscribe(ctx, pubsub.ChatChannel)
	ch := redisPubsub.Channel()
	defer redisPubsub.Close()

	for {
		select {
		case message := <-ch:
			if message.Payload == "ping" {
				s.redisClient.Publish(ctx, pubsub.ChatChannel, "pong")
			} else {
				var bMessage broadcastMessage

				err := json.Unmarshal([]byte(message.Payload), &bMessage)
				if err != nil {
					log.Println(err)
					continue
				}

				userClients := s.userClientsMap[bMessage.RecipientID]

				for _, userClient := range userClients {
					userClient.SendNotification(&bMessage.Notification)
				}
			}

		case client := <-s.register:
			userClients := s.userClientsMap[client.userID]

			if userClients == nil {
				userClients = make(map[uuid.UUID]*Client)
				s.userClientsMap[client.userID] = userClients
			}
			userClients[client.Id] = client

		case client := <-s.unregister:
			if _, ok := s.userClientsMap[client.userID]; ok {
				delete(s.userClientsMap[client.userID], client.Id)
				close(client.sendChannel)
			}
		}

	}
}

func (s *hub) BroadcastToClients(notification OutgoingNotification, recipientIDs []uuid.UUID) {
	for _, recipientID := range recipientIDs {
		message := broadcastMessage{
			Notification: notification,
			RecipientID:  recipientID,
		}

		json, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		s.redisClient.Publish(ctx, pubsub.ChatChannel, []byte(json))
	}
}

func (s *hub) RegisterClient(client *Client) {
	s.register <- client
}

func (s *hub) UnregisterClient(client *Client) {
	s.unregister <- client
}
