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

func (h *Room) Run() {
	for {
		select {
		case client := <-h.Join:
			h.Clients[client] = true

			message := domain.NewMessage(fmt.Sprintf("%s %s joined", client.Sender.Avatar, client.Sender.Name), "system", client.Sender)

			for currentClient := range h.Clients {
				if currentClient.Id == client.Id {
					continue
				}
				select {
				case currentClient.Send <- message:
				default:
					close(currentClient.Send)
					delete(h.Clients, currentClient)
				}
			}
		case client := <-h.Leave:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

			message := domain.NewMessage(fmt.Sprintf("%s %s left", client.Sender.Avatar, client.Sender.Name), "system", client.Sender)

			for currentClient := range h.Clients {
				select {
				case currentClient.Send <- message:
				default:
					close(currentClient.Send)
					delete(h.Clients, currentClient)
				}
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
