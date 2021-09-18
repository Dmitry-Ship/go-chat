package application

import (
	"GitHub/go-chat/server/domain"
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

func (room *Room) Run() {
	for {
		select {
		case client := <-room.Join:
			room.Clients[client] = true

			message := domain.NewMessage(fmt.Sprintf("%s %s joined", client.Sender.Avatar, client.Sender.Name), "system", client.Sender)

			for currentClient := range room.Clients {
				if currentClient.Id == client.Id {
					continue
				}
				select {
				case currentClient.Send <- message:
				default:
					close(currentClient.Send)
					delete(room.Clients, currentClient)
				}
			}
		case client := <-room.Leave:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				close(client.Send)
			}

			message := domain.NewMessage(fmt.Sprintf("%s %s left", client.Sender.Avatar, client.Sender.Name), "system", client.Sender)

			for currentClient := range room.Clients {
				select {
				case currentClient.Send <- message:
				default:
					close(currentClient.Send)
					delete(room.Clients, currentClient)
				}
			}
		case message := <-room.Broadcast:
			for client := range room.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(room.Clients, client)
				}
			}
		}
	}
}
