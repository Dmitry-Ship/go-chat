package ws

import (
	"context"
	"testing"

	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockParticipantRepository struct {
	conversationIDs map[uuid.UUID][]uuid.UUID
}

func (m *mockParticipantRepository) Store(ctx context.Context, participant *domain.Participant) error {
	return nil
}

func (m *mockParticipantRepository) Delete(ctx context.Context, participantID uuid.UUID) error {
	return nil
}

func (m *mockParticipantRepository) GetByConversationIDAndUserID(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	return nil, nil
}

func (m *mockParticipantRepository) GetIDsByConversationID(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

func (m *mockParticipantRepository) GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return m.conversationIDs[userID], nil
}

func TestNewActiveClients(t *testing.T) {
	ac := NewActiveClients(nil)

	assert.NotNil(t, ac)
	assert.IsType(t, &activeClients{}, ac)
}

func TestActiveClients_AddClient(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	client := &Client{
		Id:     uuid.New(),
		UserID: uuid.New(),
	}

	clientID := ac.AddClient(client)

	assert.Equal(t, client.Id, clientID)
	assert.NotNil(t, ac.GetClient(client.Id))
}

func TestActiveClients_GetClient(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	client := &Client{
		Id:     uuid.New(),
		UserID: uuid.New(),
	}

	ac.AddClient(client)

	retrieved := ac.GetClient(client.Id)
	assert.NotNil(t, retrieved)
	assert.Equal(t, client.Id, retrieved.Id)
}

func TestActiveClients_GetClient_NotFound(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	retrieved := ac.GetClient(uuid.New())

	assert.Nil(t, retrieved)
}

func TestActiveClients_RemoveClient(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	client := &Client{
		Id:          uuid.New(),
		UserID:      uuid.New(),
		channelIDs:  make(map[uuid.UUID]struct{}),
		sendChannel: make(chan OutgoingNotification, 1),
	}

	ac.AddClient(client)
	ac.RemoveClient(client)

	assert.Nil(t, ac.GetClient(client.Id))
}

func TestActiveClients_AddClientToChannel(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	client := &Client{
		Id:         uuid.New(),
		UserID:     uuid.New(),
		channelIDs: make(map[uuid.UUID]struct{}),
	}
	channelID := uuid.New()

	ac.AddClient(client)
	ac.AddClientToChannel(client, channelID)

	clients := ac.GetClientsByChannel(channelID)
	assert.Len(t, clients, 1)
	assert.Equal(t, client.Id, clients[0].Id)
}

func TestActiveClients_RemoveClientFromChannel(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	client := &Client{
		Id:         uuid.New(),
		UserID:     uuid.New(),
		channelIDs: make(map[uuid.UUID]struct{}),
	}
	channelID := uuid.New()

	ac.AddClient(client)
	ac.AddClientToChannel(client, channelID)
	ac.RemoveClientFromChannel(client, channelID)

	clients := ac.GetClientsByChannel(channelID)
	assert.Nil(t, clients)
}

func TestActiveClients_GetClientsByChannel_NotFound(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	clients := ac.GetClientsByChannel(uuid.New())

	assert.Nil(t, clients)
}

func TestActiveClients_GetClientsByUser(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	userID := uuid.New()
	client1 := &Client{
		Id:         uuid.New(),
		UserID:     userID,
		channelIDs: make(map[uuid.UUID]struct{}),
	}
	client2 := &Client{
		Id:         uuid.New(),
		UserID:     userID,
		channelIDs: make(map[uuid.UUID]struct{}),
	}

	ac.AddClient(client1)
	ac.AddClient(client2)

	clients := ac.GetClientsByUser(userID)
	assert.Len(t, clients, 2)
}

func TestActiveClients_GetClientsByUser_NotFound(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	clients := ac.GetClientsByUser(uuid.New())

	assert.Nil(t, clients)
}

func TestActiveClients_InvalidateMembership(t *testing.T) {
	userID := uuid.New()
	conversationID1 := uuid.New()
	conversationID2 := uuid.New()

	mockRepo := &mockParticipantRepository{
		conversationIDs: map[uuid.UUID][]uuid.UUID{
			userID: {conversationID1, conversationID2},
		},
	}
	ac := NewActiveClients(mockRepo)

	client := &Client{
		Id:          uuid.New(),
		UserID:      userID,
		channelIDs:  make(map[uuid.UUID]struct{}),
		sendChannel: make(chan OutgoingNotification, 1),
	}

	ac.AddClient(client)
	ac.AddClientToChannel(client, conversationID1)

	err := ac.InvalidateMembership(context.Background(), userID)

	assert.NoError(t, err)

	clients := ac.GetClientsByChannel(conversationID1)
	assert.Len(t, clients, 1)
	assert.Equal(t, client.Id, clients[0].Id)

	clients = ac.GetClientsByChannel(conversationID2)
	assert.Len(t, clients, 1)
	assert.Equal(t, client.Id, clients[0].Id)
}

func TestActiveClients_InvalidateMembership_NoClients(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	err := ac.InvalidateMembership(context.Background(), uuid.New())

	assert.NoError(t, err)
}

func TestActiveClients_MultipleClientsSameChannel(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	channelID := uuid.New()
	client1 := &Client{
		Id:         uuid.New(),
		UserID:     uuid.New(),
		channelIDs: make(map[uuid.UUID]struct{}),
	}
	client2 := &Client{
		Id:         uuid.New(),
		UserID:     uuid.New(),
		channelIDs: make(map[uuid.UUID]struct{}),
	}

	ac.AddClient(client1)
	ac.AddClient(client2)
	ac.AddClientToChannel(client1, channelID)
	ac.AddClientToChannel(client2, channelID)

	clients := ac.GetClientsByChannel(channelID)
	assert.Len(t, clients, 2)
}

func TestActiveClients_RemoveClient_RemovesFromUserList(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(mockRepo)

	userID := uuid.New()
	client := &Client{
		Id:          uuid.New(),
		UserID:      userID,
		channelIDs:  make(map[uuid.UUID]struct{}),
		sendChannel: make(chan OutgoingNotification, 1),
	}

	ac.AddClient(client)

	clients := ac.GetClientsByUser(userID)
	assert.Len(t, clients, 1)

	ac.RemoveClient(client)

	clients = ac.GetClientsByUser(userID)
	assert.Nil(t, clients)
}
