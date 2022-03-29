package ws

import (
	"github.com/google/uuid"
)

type Hub interface {
	BroadcastToClients(notification OutgoingNotification, recipientID uuid.UUID)
	UnregisterClient(client *Client)
	RegisterClient(client *Client)
}

type broadcastMessage struct {
	notification OutgoingNotification
	recipientID  uuid.UUID
}

type hub struct {
	broadcast      chan broadcastMessage
	register       chan *Client
	unregister     chan *Client
	userClientsMap map[uuid.UUID]map[uuid.UUID]*Client
}

func NewHub() *hub {
	return &hub{
		broadcast:      make(chan broadcastMessage, 1024),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		userClientsMap: make(map[uuid.UUID]map[uuid.UUID]*Client),
	}
}

func (s *hub) Run() {
	for {
		select {
		case broadcastMessage := <-s.broadcast:
			userClients := s.userClientsMap[broadcastMessage.recipientID]

			for _, userClient := range userClients {
				userClient.SendNotification(&broadcastMessage.notification)
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
				close(client.send)
			}
		}

	}
}

func (s *hub) BroadcastToClients(notification OutgoingNotification, recipientID uuid.UUID) {
	s.broadcast <- broadcastMessage{
		notification: notification,
		recipientID:  recipientID,
	}
}

func (s *hub) RegisterClient(client *Client) {
	s.register <- client
}

func (s *hub) UnregisterClient(client *Client) {
	s.unregister <- client
}
