package ws

import (
	"sync"

	"github.com/google/uuid"
)

type ActiveClients interface {
	AddClient(c *client)
	RemoveClient(c *client)
	SendToUserClients(userID uuid.UUID, notification OutgoingNotification)
}

type activeClients struct {
	mu             sync.RWMutex
	userClientsMap map[uuid.UUID]map[uuid.UUID]*client
}

func NewActiveClients() *activeClients {
	return &activeClients{
		userClientsMap: make(map[uuid.UUID]map[uuid.UUID]*client),
		mu:             sync.RWMutex{},
	}
}

func (ac *activeClients) AddClient(c *client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	userClients, ok := ac.userClientsMap[c.UserID]

	if !ok {
		userClients = make(map[uuid.UUID]*client)
		ac.userClientsMap[c.UserID] = userClients
	}

	userClients[c.Id] = c
}

func (ac *activeClients) RemoveClient(c *client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, ok := ac.userClientsMap[c.UserID]; ok {
		close(c.sendChannel)
		delete(ac.userClientsMap[c.UserID], c.Id)
	}
}

func (ac *activeClients) SendToUserClients(userID uuid.UUID, notification OutgoingNotification) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	userClients, ok := ac.userClientsMap[userID]

	if !ok {
		return
	}

	for _, client := range userClients {
		client.sendNotification(notification)
	}
}
