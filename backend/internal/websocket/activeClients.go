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
	RemoveClientsByUser(userID uuid.UUID)
	AddClientToChannel(c *Client, channelID uuid.UUID)
	RemoveClientFromChannel(c *Client, channelID uuid.UUID)
	GetClientsByChannel(channelID uuid.UUID) []*Client
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
}

type activeClients struct {
	ctx          context.Context
	mu           sync.RWMutex
	clients      map[clientKey]*Client
	byUserID     map[uuid.UUID]map[*Client]struct{}
	byChannelID  map[uuid.UUID]map[*Client]struct{}
	byClientID   map[uuid.UUID]*Client
	participants domain.ParticipantRepository
}

func NewActiveClients(ctx context.Context, participants domain.ParticipantRepository) *activeClients {
	return &activeClients{
		ctx:          ctx,
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

	for channelID, clients := range ac.byChannelID {
		delete(clients, c)
		if len(clients) == 0 {
			delete(ac.byChannelID, channelID)
		}
	}
}

func (ac *activeClients) RemoveClientsByUser(userID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	userClients, exists := ac.byUserID[userID]
	if !exists {
		return
	}

	for client := range userClients {
		key := clientKey{userID: userID, clientID: client.Id}
		close(client.sendChannel)
		delete(ac.clients, key)
		delete(ac.byClientID, client.Id)

		for channelID, channelClients := range ac.byChannelID {
			delete(channelClients, client)
			if len(channelClients) == 0 {
				delete(ac.byChannelID, channelID)
			}
		}
	}

	delete(ac.byUserID, userID)
}

func (ac *activeClients) AddClientToChannel(c *Client, channelID uuid.UUID) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.byChannelID[channelID]; !exists {
		ac.byChannelID[channelID] = make(map[*Client]struct{})
	}
	ac.byChannelID[channelID][c] = struct{}{}
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
	ac.mu.RLock()
	userClients, exists := ac.byUserID[userID]
	if !exists {
		ac.mu.RUnlock()
		return nil
	}

	clients := make([]*Client, 0, len(userClients))
	for client := range userClients {
		clients = append(clients, client)
	}
	ac.mu.RUnlock()

	conversationIDs, err := ac.participants.GetConversationIDsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()

	currentUserClients, stillExists := ac.byUserID[userID]
	if !stillExists {
		return nil
	}

	validClients := make(map[*Client]struct{})
	for _, client := range clients {
		if _, ok := currentUserClients[client]; ok {
			validClients[client] = struct{}{}
		}
	}

	for client := range validClients {
		var channelIDs []uuid.UUID
		for channelID, channelClients := range ac.byChannelID {
			if _, exists := channelClients[client]; exists {
				channelIDs = append(channelIDs, channelID)
			}
		}
		for _, channelID := range channelIDs {
			if channelClients, exists := ac.byChannelID[channelID]; exists {
				delete(channelClients, client)
				if len(channelClients) == 0 {
					delete(ac.byChannelID, channelID)
				}
			}
		}
	}

	for client := range validClients {
		for _, conversationID := range conversationIDs {
			if _, exists := ac.byChannelID[conversationID]; !exists {
				ac.byChannelID[conversationID] = make(map[*Client]struct{})
			}
			ac.byChannelID[conversationID][client] = struct{}{}
		}
	}

	return nil
}
