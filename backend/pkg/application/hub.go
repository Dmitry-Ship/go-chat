package application

import (
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
	broadcast  chan *notification
	Register   chan *Client
	Unregister chan *Client
	clients    map[uuid.UUID]map[uuid.UUID]*Client
}

func NewHub() *hub {
	return &hub{
		broadcast:  make(chan *notification, 1024),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]map[uuid.UUID]*Client),
	}
}

func (s *hub) Run() {
	for {
		select {
		case message := <-s.broadcast:
			clients := s.clients[message.UserId]

			for _, client := range clients {
				client.SendNotification(message.Type, message.Payload)
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
	s.broadcast <- &notification{
		Type:    notificationType,
		Payload: payload,
		UserId:  userID,
	}
}
