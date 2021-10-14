package application

import (
	"github.com/google/uuid"
)

type message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	UserId  uuid.UUID   `json:"userId"`
}

type Hub interface {
	Register(client *Client)
	Unregister(client *Client)
	BroadcastMessage(messageType string, payload interface{}, userID uuid.UUID)
	Run()
}

type hub struct {
	broadcast  chan *message
	register   chan *Client
	unregister chan *Client
	clients    map[uuid.UUID][]*Client
}

func NewHub() *hub {
	return &hub{
		broadcast:  make(chan *message, 1024),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID][]*Client),
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

		case client := <-s.register:
			s.clients[client.userID] = append(s.clients[client.userID], client)

		case client := <-s.unregister:
			for i, c := range s.clients[client.userID] {
				if c == client {
					s.clients[client.userID] = append(s.clients[client.userID][:i], s.clients[client.userID][i+1:]...)
					break
				}
			}

		}

	}
}

func (s *hub) Register(client *Client) {
	s.register <- client
}

func (s *hub) Unregister(client *Client) {
	s.unregister <- client
}

func (s *hub) BroadcastMessage(messageType string, payload interface{}, userID uuid.UUID) {
	s.broadcast <- &message{
		Type:    messageType,
		Payload: payload,
		UserId:  userID,
	}
}
