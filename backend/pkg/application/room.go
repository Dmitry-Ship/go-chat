package application

import (
	"GitHub/go-chat/backend/domain"
	"fmt"
)

// Room maintains the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	Clients   map[*Client]bool
	Broadcast chan domain.Message
	Join      chan *Client
	Leave     chan *Client
}

func NewRoom() *Room {
	return &Room{
		Broadcast: make(chan domain.Message, 1024),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Clients:   make(map[*Client]bool),
	}
}

func (room *Room) broadcastMessage(message interface{}) {
	for currentClient := range room.Clients {
		notification := currentClient.NewNotification("message", message)

		select {
		case currentClient.Send <- notification:
		default:
			close(currentClient.Send)
			delete(room.Clients, currentClient)
		}
	}
}

func (room *Room) Run() {
	for {
		select {
		case client := <-room.Join:
			room.Clients[client] = true

			message := domain.NewMessage(fmt.Sprintf("%s %s joined", client.Sender.Avatar, client.Sender.Name), "system", client.Sender)

			room.broadcastMessage(message)
		case client := <-room.Leave:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				close(client.Send)
			}

			message := domain.NewMessage(fmt.Sprintf("%s %s left", client.Sender.Avatar, client.Sender.Name), "system", client.Sender)

			room.broadcastMessage(message)
		case message := <-room.Broadcast:
			room.broadcastMessage(message)
		}
	}
}
