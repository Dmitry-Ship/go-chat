package ws

import (
	"context"
	"encoding/json"
	"log"

	pubsub "GitHub/go-chat/backend/internal/infra/redis"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type ConnectionsPool interface {
	RegisterClient(conn *websocket.Conn, userID uuid.UUID) error
	BroadcastToUsers(userIDs []uuid.UUID, notification OutgoingNotification) error
}

type broadcastMessage struct {
	Payload OutgoingNotification `json:"notification"`
	UserIDs []uuid.UUID          `json:"user_ids"`
}

type connectionsPool struct {
	ctx                   context.Context
	registerClient        chan *client
	unregisterClient      chan *client
	IncomingNotifications chan *IncomingNotification
	userClientsMap        map[uuid.UUID]map[uuid.UUID]*client
	redisClient           *redis.Client
}

func NewConnectionsPool(ctx context.Context, redisClient *redis.Client) *connectionsPool {
	return &connectionsPool{
		ctx:                   ctx,
		registerClient:        make(chan *client, 1024),
		unregisterClient:      make(chan *client, 1024),
		userClientsMap:        make(map[uuid.UUID]map[uuid.UUID]*client),
		redisClient:           redisClient,
		IncomingNotifications: make(chan *IncomingNotification, 1024),
	}
}

func (s *connectionsPool) Run() {
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

			var bMessage broadcastMessage

			err := json.Unmarshal([]byte(message.Payload), &bMessage)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, userID := range bMessage.UserIDs {
				clients, ok := s.userClientsMap[userID]

				if !ok {
					continue
				}

				for _, client := range clients {
					client.SendNotification(&bMessage.Payload)
				}
			}

		case newClient := <-s.registerClient:
			userClients, ok := s.userClientsMap[newClient.UserID]

			if !ok {
				userClients = make(map[uuid.UUID]*client)
				s.userClientsMap[newClient.UserID] = userClients
			}

			userClients[newClient.Id] = newClient

		case client := <-s.unregisterClient:
			if _, ok := s.userClientsMap[client.UserID]; ok {
				delete(s.userClientsMap[client.UserID], client.Id)
				close(client.sendChannel)
			}
		}
	}
}

func (s *connectionsPool) BroadcastToUsers(userIDs []uuid.UUID, notification OutgoingNotification) error {
	message := broadcastMessage{
		Payload: notification,
		UserIDs: userIDs,
	}

	json, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return s.redisClient.Publish(s.ctx, pubsub.ChatChannel, []byte(json)).Err()
}

func (s *connectionsPool) RegisterClient(conn *websocket.Conn, userID uuid.UUID) error {
	client := NewClient(conn, s.unregisterClient, s.IncomingNotifications, userID)

	go client.WritePump()
	go client.ReadPump()

	s.registerClient <- client

	return nil
}
