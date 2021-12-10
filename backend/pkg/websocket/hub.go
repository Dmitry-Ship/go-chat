package ws

import (
	"GitHub/go-chat/backend/pkg/redis"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type notification struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	UserId  uuid.UUID   `json:"userId"`
}

type HubBroadcaster interface {
	BroadcastNotification(notificationType string, payload interface{}, userID uuid.UUID)
}

type hub struct {
	Register    chan *Client
	Unregister  chan *Client
	clients     map[uuid.UUID]map[uuid.UUID]*Client
	redisClient redis.RedisClient
}

func NewHub(redisClient redis.RedisClient) *hub {
	return &hub{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		clients:     make(map[uuid.UUID]map[uuid.UUID]*Client),
		redisClient: redisClient,
	}
}

func (s *hub) Run() {
	redisMessageChannel := s.redisClient.GetMessageChannel()
	for {
		select {
		case message := <-redisMessageChannel:
			var n notification
			err := json.Unmarshal([]byte(message), &n)
			if err != nil {
				log.Println(err)
			}

			clients := s.clients[n.UserId]

			for _, client := range clients {
				client.SendNotification(n.Type, n.Payload)
			}

		case client := <-s.Register:
			userClients := s.clients[client.userID]
			if userClients == nil {
				userClients = make(map[uuid.UUID]*Client)
				s.clients[client.userID] = userClients
			}
			userClients[client.Id] = client

		case client := <-s.Unregister:
			if _, ok := s.clients[client.userID]; ok {
				delete(s.clients[client.userID], client.Id)
				close(client.send)
			}
		}

	}
}

func (s *hub) BroadcastNotification(notificationType string, payload interface{}, userID uuid.UUID) {
	n := &notification{
		Type:    notificationType,
		Payload: payload,
		UserId:  userID,
	}

	s.redisClient.SendToChannel("chat", n)
}
