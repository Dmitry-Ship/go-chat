package ws

import (
	"sync"

	"github.com/google/uuid"
)

type clientKey struct {
	userID   uuid.UUID
	clientID uuid.UUID
}

type ActiveClients interface {
	AddClient(c *client)
	RemoveClient(c *client)
	SendToUserClients(userID uuid.UUID, notification OutgoingNotification)
}

type activeClients struct {
	mu      sync.RWMutex
	clients map[clientKey]*client
}

func NewActiveClients() *activeClients {
	return &activeClients{
		clients: make(map[clientKey]*client),
		mu:      sync.RWMutex{},
	}
}

func (ac *activeClients) AddClient(c *client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	key := clientKey{userID: c.UserID, clientID: c.Id}
	ac.clients[key] = c
}

func (ac *activeClients) RemoveClient(c *client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	key := clientKey{userID: c.UserID, clientID: c.Id}
	close(c.sendChannel)
	delete(ac.clients, key)
}

func (ac *activeClients) SendToUserClients(userID uuid.UUID, notification OutgoingNotification) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	for key, client := range ac.clients {
		if key.userID == userID {
			client.sendNotification(notification)
		}
	}
}
