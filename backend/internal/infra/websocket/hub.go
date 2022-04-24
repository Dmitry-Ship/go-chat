package ws

import (
	"context"
	"encoding/json"
	"log"

	pubsub "GitHub/go-chat/backend/internal/infra/redis"

	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
)

type Hub interface {
	BroadcastToTopic(topic string, notification OutgoingNotification)
	SubscribeToTopic(topic string, userID uuid.UUID)
	UnsubscribeFromTopic(topic string, userID uuid.UUID)
	UnregisterClient(client *Client)
	RegisterClient(client *Client)
}

type broadcastMessageToTopic struct {
	Notification OutgoingNotification `json:"notification"`
	Topic        string               `json:"topic"`
}

type subscrition struct {
	userID uuid.UUID
	Topic  string
}

type hub struct {
	register         chan *Client
	unregister       chan *Client
	subscribe        chan *subscrition
	unsubscribe      chan *subscrition
	clientsTopicsMap map[string]map[uuid.UUID]*Client
	userClientsMap   map[uuid.UUID]map[uuid.UUID]*Client
	redisClient      *redis.Client
}

var ctx = context.Background()

func NewHub(redisClient *redis.Client) *hub {
	return &hub{
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		subscribe:        make(chan *subscrition),
		unsubscribe:      make(chan *subscrition),
		clientsTopicsMap: make(map[string]map[uuid.UUID]*Client),
		userClientsMap:   make(map[uuid.UUID]map[uuid.UUID]*Client),
		redisClient:      redisClient,
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
				continue
			}

			var bMessage broadcastMessageToTopic

			err := json.Unmarshal([]byte(message.Payload), &bMessage)

			if err != nil {
				log.Println(err)
				continue
			}

			topicClients := s.clientsTopicsMap[bMessage.Topic]

			for _, client := range topicClients {
				client.SendNotification(&bMessage.Notification)
			}

		case client := <-s.register:
			userClients, ok := s.userClientsMap[client.UserID]

			if !ok {
				userClients = make(map[uuid.UUID]*Client)
				s.userClientsMap[client.UserID] = userClients
			}

			userClients[client.Id] = client

		case client := <-s.unregister:
			if _, ok := s.userClientsMap[client.UserID]; ok {
				delete(s.userClientsMap[client.UserID], client.Id)
				close(client.sendChannel)
			}

		case subscrition := <-s.subscribe:
			clientsMap, ok := s.userClientsMap[subscrition.userID]

			if !ok {
				return
			}

			if _, ok := s.clientsTopicsMap[subscrition.Topic]; !ok {
				s.clientsTopicsMap[subscrition.Topic] = make(map[uuid.UUID]*Client)
			}

			for _, client := range clientsMap {
				s.clientsTopicsMap[subscrition.Topic][client.Id] = client
			}

		case subscrition := <-s.unsubscribe:
			clientsMap, ok := s.userClientsMap[subscrition.userID]

			if !ok {
				return
			}

			if _, ok := s.clientsTopicsMap[subscrition.Topic]; ok {
				for _, client := range clientsMap {
					delete(s.clientsTopicsMap[subscrition.Topic], client.Id)
				}
			}
		}
	}
}

func (s *hub) BroadcastToTopic(topic string, notification OutgoingNotification) {
	message := broadcastMessageToTopic{
		Notification: notification,
		Topic:        topic,
	}

	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	s.redisClient.Publish(ctx, pubsub.ChatChannel, []byte(json))
}

func (s *hub) SubscribeToTopic(topic string, userID uuid.UUID) {
	s.subscribe <- &subscrition{
		userID: userID,
		Topic:  topic,
	}
}

func (s *hub) UnsubscribeFromTopic(topic string, userID uuid.UUID) {
	s.unsubscribe <- &subscrition{
		userID: userID,
		Topic:  topic,
	}
}

func (s *hub) RegisterClient(client *Client) {
	s.register <- client
}

func (s *hub) UnregisterClient(client *Client) {
	s.unregister <- client
}
