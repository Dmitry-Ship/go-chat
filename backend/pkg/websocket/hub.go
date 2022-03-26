package ws

import (
	"github.com/google/uuid"
)

type HubBroadcaster interface {
	BroadcastNotification(notification OutgoingNotification)
}

type hub struct {
	broadcast  chan *OutgoingNotification
	Register   chan *Client
	Unregister chan *Client
	clients    map[uuid.UUID]map[uuid.UUID]*Client
}

func NewHub() *hub {
	return &hub{
		broadcast:  make(chan *OutgoingNotification, 1024),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]map[uuid.UUID]*Client),
	}
}

func (s *hub) Run() {
	for {
		select {
		case notification := <-s.broadcast:
			clients := s.clients[notification.UserID]

			for _, client := range clients {
				client.SendNotification(notification)
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

func (s *hub) BroadcastNotification(notification OutgoingNotification) {
	s.broadcast <- &notification
}
