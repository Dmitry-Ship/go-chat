package application

import (
	"GitHub/go-chat/server/domain"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Clients   map[*Client]bool
	Broadcast chan domain.Message
	Join      chan *Client
	Leave     chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast: make(chan domain.Message, 1024),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Clients:   make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Join:
			h.Clients[client] = true
			message := domain.NewMessage("new user joined", "system", "0")

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
