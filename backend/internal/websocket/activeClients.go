package ws

import (
	"context"
	"log"
	"sync"

	"GitHub/go-chat/backend/internal/repository"

	"github.com/google/uuid"
)

type ActiveClients interface {
	AddClient(c *Client) uuid.UUID
	RemoveClient(c *Client)
	InvalidateMembership(ctx context.Context, userID uuid.UUID) error
	NotifyChannelClients(ctx context.Context, channelID uuid.UUID, notification OutgoingNotification)
}

type activeClients struct {
	mu               sync.RWMutex
	byUserID         map[uuid.UUID]map[*Client]struct{}
	byChannelID      map[uuid.UUID]map[*Client]struct{}
	byClientChannels map[*Client]map[uuid.UUID]struct{}
	participants     repository.ParticipantRepository
}

func NewActiveClients(ctx context.Context, participants repository.ParticipantRepository) *activeClients {
	return &activeClients{
		byUserID:         make(map[uuid.UUID]map[*Client]struct{}),
		byChannelID:      make(map[uuid.UUID]map[*Client]struct{}),
		byClientChannels: make(map[*Client]map[uuid.UUID]struct{}),
		participants:     participants,
	}
}

func (ac *activeClients) AddClient(c *Client) uuid.UUID {
	var conversationIDs []uuid.UUID
	if ac.participants != nil {
		ids, err := ac.participants.GetConversationIDsByUserID(context.Background(), c.UserID)
		if err == nil {
			conversationIDs = ids
		}
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()

	userClients := ac.byUserID[c.UserID]
	if userClients == nil {
		userClients = make(map[*Client]struct{})
		ac.byUserID[c.UserID] = userClients
	}
	userClients[c] = struct{}{}

	channelIDs := ac.byClientChannels[c]
	if channelIDs == nil {
		channelIDs = make(map[uuid.UUID]struct{})
		ac.byClientChannels[c] = channelIDs
	}

	for _, conversationID := range conversationIDs {
		if _, exists := ac.byChannelID[conversationID]; !exists {
			ac.byChannelID[conversationID] = make(map[*Client]struct{})
		}
		ac.byChannelID[conversationID][c] = struct{}{}
		channelIDs[conversationID] = struct{}{}
	}

	return c.Id
}

func (ac *activeClients) RemoveClient(c *Client) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	close(c.sendChannel)

	if userClients, exists := ac.byUserID[c.UserID]; exists {
		delete(userClients, c)
		if len(userClients) == 0 {
			delete(ac.byUserID, c.UserID)
		}
	}

	if channelIDs, exists := ac.byClientChannels[c]; exists {
		for channelID := range channelIDs {
			if clients, ok := ac.byChannelID[channelID]; ok {
				delete(clients, c)
				if len(clients) == 0 {
					delete(ac.byChannelID, channelID)
				}
			}
		}
		delete(ac.byClientChannels, c)
	}
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

	desiredChannels := make(map[uuid.UUID]struct{}, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		desiredChannels[conversationID] = struct{}{}
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()

	currentUserClients, stillExists := ac.byUserID[userID]
	if !stillExists {
		return nil
	}

	for _, client := range clients {
		if _, ok := currentUserClients[client]; !ok {
			continue
		}

		currentChannels := ac.byClientChannels[client]
		if currentChannels == nil {
			currentChannels = make(map[uuid.UUID]struct{})
			ac.byClientChannels[client] = currentChannels
		}

		for channelID := range currentChannels {
			if _, keep := desiredChannels[channelID]; keep {
				continue
			}
			if channelClients, exists := ac.byChannelID[channelID]; exists {
				delete(channelClients, client)
				if len(channelClients) == 0 {
					delete(ac.byChannelID, channelID)
				}
			}
			delete(currentChannels, channelID)
		}

		for channelID := range desiredChannels {
			if _, exists := currentChannels[channelID]; exists {
				continue
			}
			if _, exists := ac.byChannelID[channelID]; !exists {
				ac.byChannelID[channelID] = make(map[*Client]struct{})
			}
			ac.byChannelID[channelID][client] = struct{}{}
			currentChannels[channelID] = struct{}{}
		}
	}

	return nil
}

func (ac *activeClients) NotifyChannelClients(ctx context.Context, channelID uuid.UUID, notification OutgoingNotification) {
	ac.mu.RLock()
	clients, exists := ac.byChannelID[channelID]
	if !exists {
		ac.mu.RUnlock()
		return
	}

	defer ac.mu.RUnlock()

	for client := range clients {
		if err := client.SendNotification(notification); err != nil {
			log.Printf("Error sending notification to client %s: %v", client.Id, err)
		}
	}
}
