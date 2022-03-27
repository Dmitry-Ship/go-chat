package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type Hub interface {
	BroadcastNotification(notification OutgoingNotification)
	UnregisterClient(client *Client)
	RegisterClient(client *Client)
	SendMessage(notification IncomingNotification)
}

type WSHandler func(message IncomingNotification, data json.RawMessage)

type hub struct {
	broadcast       chan *OutgoingNotification
	messageChannel  chan IncomingNotification
	register        chan *Client
	unregister      chan *Client
	clients         map[uuid.UUID]map[uuid.UUID]*Client
	messageHandlers map[string]WSHandler
}

func NewHub() *hub {
	return &hub{
		broadcast:       make(chan *OutgoingNotification, 1024),
		messageChannel:  make(chan IncomingNotification, 100),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[uuid.UUID]map[uuid.UUID]*Client),
		messageHandlers: make(map[string]WSHandler),
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

		case client := <-s.register:
			userClients := s.clients[client.userID]

			if userClients == nil {
				userClients = make(map[uuid.UUID]*Client)
				s.clients[client.userID] = userClients
			}
			userClients[client.Id] = client

		case client := <-s.unregister:
			if _, ok := s.clients[client.userID]; ok {
				delete(s.clients[client.userID], client.Id)
				close(client.send)
			}

		case message := <-s.messageChannel:
			var data json.RawMessage

			notification := struct {
				Type string      `json:"type"`
				Data interface{} `json:"data"`
			}{
				Data: &data,
			}

			if err := json.Unmarshal(message.Data, &notification); err != nil {
				log.Println(err)
				continue
			}

			if handler, ok := s.messageHandlers[notification.Type]; ok {
				go handler(message, data)
			}
		}

	}
}

func (s *hub) BroadcastNotification(notification OutgoingNotification) {
	s.broadcast <- &notification
}

func (s *hub) RegisterClient(client *Client) {
	s.register <- client
}

func (s *hub) UnregisterClient(client *Client) {
	s.unregister <- client
}

func (s *hub) SendMessage(notification IncomingNotification) {
	s.messageChannel <- notification
}

func (s *hub) SetWSHandler(notificationType string, handler func(IncomingNotification, json.RawMessage)) {
	s.messageHandlers[notificationType] = handler
}
