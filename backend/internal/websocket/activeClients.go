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
	AddClient(c *Client) uuid.UUID
	RemoveClient(c *Client)
	GetClient(clientID uuid.UUID) *Client
	SendToUserClients(userID uuid.UUID, notification OutgoingNotification)
	SubscribeChannel(c *Client, channelID uuid.UUID)
	UnsubscribeChannel(c *Client, channelID uuid.UUID)
	UnsubscribeClientAllChannels(c *Client)
	SendToChannel(channelID uuid.UUID, notification OutgoingNotification)
	GetClientsByChannel(channelID uuid.UUID) []*Client
	GetClientsByUser(userID uuid.UUID) []*Client
	SubscribeUserToChannel(userID uuid.UUID, channelID uuid.UUID)
	UnsubscribeUserFromChannel(userID uuid.UUID, channelID uuid.UUID)
}

type activeClients struct {
	mu          sync.RWMutex
	clients     map[clientKey]*Client
	byUserID    map[uuid.UUID]map[*Client]struct{}
	byChannelID map[uuid.UUID]map[*Client]struct{}
}

func NewActiveClients() *activeClients {
	return &activeClients{
		clients:     make(map[clientKey]*Client),
		byUserID:    make(map[uuid.UUID]map[*Client]struct{}),
		byChannelID: make(map[uuid.UUID]map[*Client]struct{}),
		mu:          sync.RWMutex{},
	}
}

func (ac *activeClients) AddClient(c *Client) uuid.UUID {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	key := clientKey{userID: c.UserID, clientID: c.Id}
	ac.clients[key] = c

	if _, exists := ac.byUserID[c.UserID]; !exists {
		ac.byUserID[c.UserID] = make(map[*Client]struct{})
	}
	ac.byUserID[c.UserID][c] = struct{}{}

	return c.Id
}

func (ac *activeClients) GetClient(clientID uuid.UUID) *Client {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	for key, client := range ac.clients {
		if key.clientID == clientID {
			return client
		}
	}

	return nil
}

func (ac *activeClients) RemoveClient(c *Client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	key := clientKey{userID: c.UserID, clientID: c.Id}
	close(c.sendChannel)
	delete(ac.clients, key)

	if userClients, exists := ac.byUserID[c.UserID]; exists {
		delete(userClients, c)
		if len(userClients) == 0 {
			delete(ac.byUserID, c.UserID)
		}
	}

	ac.UnsubscribeClientAllChannels(c)
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

func (ac *activeClients) SubscribeChannel(c *Client, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.byChannelID[channelID]; !exists {
		ac.byChannelID[channelID] = make(map[*Client]struct{})
	}
	ac.byChannelID[channelID][c] = struct{}{}
}

func (ac *activeClients) UnsubscribeChannel(c *Client, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if channelClients, exists := ac.byChannelID[channelID]; exists {
		delete(channelClients, c)
		if len(channelClients) == 0 {
			delete(ac.byChannelID, channelID)
		}
	}
}

func (ac *activeClients) UnsubscribeClientAllChannels(c *Client) {
	for channelID, clients := range ac.byChannelID {
		delete(clients, c)
		if len(clients) == 0 {
			delete(ac.byChannelID, channelID)
		}
	}
}

func (ac *activeClients) SendToChannel(channelID uuid.UUID, notification OutgoingNotification) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	channelClients, exists := ac.byChannelID[channelID]
	if !exists {
		return
	}

	for client := range channelClients {
		client.sendNotification(notification)
	}
}

func (ac *activeClients) GetClientsByChannel(channelID uuid.UUID) []*Client {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	channelClients, exists := ac.byChannelID[channelID]
	if !exists {
		return nil
	}

	clients := make([]*Client, 0, len(channelClients))
	for client := range channelClients {
		clients = append(clients, client)
	}

	return clients
}

func (ac *activeClients) GetClientsByUser(userID uuid.UUID) []*Client {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	userClients, exists := ac.byUserID[userID]
	if !exists {
		return nil
	}

	clients := make([]*Client, 0, len(userClients))
	for client := range userClients {
		clients = append(clients, client)
	}

	return clients
}

func (ac *activeClients) SubscribeUserToChannel(userID uuid.UUID, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	userClients, exists := ac.byUserID[userID]
	if !exists {
		return
	}

	if _, exists := ac.byChannelID[channelID]; !exists {
		ac.byChannelID[channelID] = make(map[*Client]struct{})
	}

	for client := range userClients {
		ac.byChannelID[channelID][client] = struct{}{}
	}
}

func (ac *activeClients) UnsubscribeUserFromChannel(userID uuid.UUID, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	userClients, exists := ac.byUserID[userID]
	if !exists {
		return
	}

	channelClients, exists := ac.byChannelID[channelID]
	if !exists {
		return
	}

	for client := range userClients {
		delete(channelClients, client)
	}

	if len(channelClients) == 0 {
		delete(ac.byChannelID, channelID)
	}
}
