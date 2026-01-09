package ws

import (
	"context"
	"sync"

	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type clientKey struct {
	userID   uuid.UUID
	clientID uuid.UUID
}

type ActiveClients interface {
	AddClient(c *Client) uuid.UUID
	RemoveClient(c *Client)
	AddClientToChannel(c *Client, channelID uuid.UUID)
	RemoveClientFromChannel(c *Client, channelID uuid.UUID)
	GetClient(clientID uuid.UUID) *Client
	GetClientsByChannel(channelID uuid.UUID) []*Client
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
	GetClientsByUser(userID uuid.UUID) []*Client
}

type activeClients struct {
	mu           sync.RWMutex
	clients      map[clientKey]*Client
	byUserID     map[uuid.UUID]map[*Client]struct{}
	byChannelID  map[uuid.UUID]map[*Client]struct{}
	byClientID   map[uuid.UUID]*Client
	participants domain.ParticipantRepository
}

func NewActiveClients(participants domain.ParticipantRepository) *activeClients {
	return &activeClients{
		clients:      make(map[clientKey]*Client),
		byUserID:     make(map[uuid.UUID]map[*Client]struct{}),
		byChannelID:  make(map[uuid.UUID]map[*Client]struct{}),
		byClientID:   make(map[uuid.UUID]*Client),
		participants: participants,
	}
}

func (ac *activeClients) AddClient(c *Client) uuid.UUID {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	key := clientKey{userID: c.UserID, clientID: c.Id}
	ac.clients[key] = c
	ac.byClientID[c.Id] = c

	userClients := ac.byUserID[c.UserID]
	if userClients == nil {
		userClients = make(map[*Client]struct{})
		ac.byUserID[c.UserID] = userClients
	}
	userClients[c] = struct{}{}

	return c.Id
}

func (ac *activeClients) GetClient(clientID uuid.UUID) *Client {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	return ac.byClientID[clientID]
}

func (ac *activeClients) RemoveClient(c *Client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	key := clientKey{userID: c.UserID, clientID: c.Id}
	close(c.sendChannel)
	delete(ac.clients, key)
	delete(ac.byClientID, c.Id)

	if userClients, exists := ac.byUserID[c.UserID]; exists {
		delete(userClients, c)
		if len(userClients) == 0 {
			delete(ac.byUserID, c.UserID)
		}
	}

	for channelID := range c.channelIDs {
		if clients, exists := ac.byChannelID[channelID]; exists {
			delete(clients, c)
			if len(clients) == 0 {
				delete(ac.byChannelID, channelID)
			}
		}
	}
}

func (ac *activeClients) AddClientToChannel(c *Client, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.byChannelID[channelID]; !exists {
		ac.byChannelID[channelID] = make(map[*Client]struct{})
	}
	ac.byChannelID[channelID][c] = struct{}{}
	c.channelIDs[channelID] = struct{}{}
}

func (ac *activeClients) RemoveClientFromChannel(c *Client, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if clients, exists := ac.byChannelID[channelID]; exists {
		delete(clients, c)
		if len(clients) == 0 {
			delete(ac.byChannelID, channelID)
		}
	}
	delete(c.channelIDs, channelID)
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

func (ac *activeClients) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	ac.mu.Lock()
	userClients, exists := ac.byUserID[userID]
	if !exists {
		ac.mu.Unlock()
		return nil
	}

	clients := make([]*Client, 0, len(userClients))
	for client := range userClients {
		clients = append(clients, client)
	}
	ac.mu.Unlock()

	conversationIDs, err := ac.participants.GetConversationIDsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()

	for _, client := range clients {
		for channelID := range client.channelIDs {
			if channelClients, exists := ac.byChannelID[channelID]; exists {
				delete(channelClients, client)
				if len(channelClients) == 0 {
					delete(ac.byChannelID, channelID)
				}
			}
		}
		client.channelIDs = make(map[uuid.UUID]struct{})
	}

	for _, client := range clients {
		for _, conversationID := range conversationIDs {
			if _, exists := ac.byChannelID[conversationID]; !exists {
				ac.byChannelID[conversationID] = make(map[*Client]struct{})
			}
			ac.byChannelID[conversationID][client] = struct{}{}
			client.channelIDs[conversationID] = struct{}{}
		}
	}

	return nil
}
