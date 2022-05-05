package ws

import (
	"sync"

	"github.com/google/uuid"
)

type activeClients struct {
	mu             sync.RWMutex
	userClientsMap map[uuid.UUID]map[uuid.UUID]*client
}

func newActiveClients() *activeClients {
	return &activeClients{
		userClientsMap: make(map[uuid.UUID]map[uuid.UUID]*client),
		mu:             sync.RWMutex{},
	}
}

func (ac *activeClients) addClient(c *client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	userClients, ok := ac.userClientsMap[c.UserID]

	if !ok {
		userClients = make(map[uuid.UUID]*client)
		ac.userClientsMap[c.UserID] = userClients
	}

	userClients[c.Id] = c
}

func (ac *activeClients) removeClient(c *client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, ok := ac.userClientsMap[c.UserID]; ok {
		delete(ac.userClientsMap[c.UserID], c.Id)
		close(c.sendChannel)
	}
}

func (ac *activeClients) sendToUserClients(userID uuid.UUID, notification OutgoingNotification) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	userClients, ok := ac.userClientsMap[userID]

	if !ok {
		return
	}

	for _, client := range userClients {
		client.SendNotification(&notification)
	}
}
