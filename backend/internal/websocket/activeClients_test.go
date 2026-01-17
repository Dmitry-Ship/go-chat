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
	ac := NewActiveClients(context.Background(), nil)

	assert.NotNil(t, ac)
	assert.IsType(t, &activeClients{}, ac)
}

func TestActiveClients_AddClient(t *testing.T) {
	channelID := uuid.New()
	userID := uuid.New()
	mockRepo := &mockParticipantRepository{
		conversationIDs: map[uuid.UUID][]uuid.UUID{
			userID: {channelID},
		},
	}
	ac := NewActiveClients(context.Background(), mockRepo)

	client := &Client{
		Id:     uuid.New(),
		UserID: userID,
	}

	clientID := ac.AddClient(client)

	assert.Equal(t, client.Id, clientID)

	clients := getClientsByChannel(ac, channelID)
	assert.Len(t, clients, 1)
	assert.Equal(t, client.Id, clients[0].Id)
}

func TestActiveClients_RemoveClient(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(context.Background(), mockRepo)

	client := &Client{
		Id:     uuid.New(),
		UserID: uuid.New(),

		sendChannel: make(chan OutgoingNotification, 1),
	}

	ac.AddClient(client)
	ac.RemoveClient(client)
}

func getClientsByChannel(ac *activeClients, channelID uuid.UUID) []*Client {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	channelClients := ac.byChannelID[channelID]
	if len(channelClients) == 0 {
		return nil
	}

	clients := make([]*Client, 0, len(channelClients))
	for client := range channelClients {
		clients = append(clients, client)
	}

	return clients
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
	ac := NewActiveClients(context.Background(), mockRepo)

	client := &Client{
		Id:     uuid.New(),
		UserID: userID,

		sendChannel: make(chan OutgoingNotification, 1),
	}

	ac.AddClient(client)

	err := ac.InvalidateMembership(context.Background(), userID)

	assert.NoError(t, err)

	clients := getClientsByChannel(ac, conversationID1)
	assert.Len(t, clients, 1)
	assert.Equal(t, client.Id, clients[0].Id)

	clients = getClientsByChannel(ac, conversationID2)
	assert.Len(t, clients, 1)
	assert.Equal(t, client.Id, clients[0].Id)
}

func TestActiveClients_InvalidateMembership_NoClients(t *testing.T) {
	mockRepo := &mockParticipantRepository{
		conversationIDs: make(map[uuid.UUID][]uuid.UUID),
	}
	ac := NewActiveClients(context.Background(), mockRepo)

	err := ac.InvalidateMembership(context.Background(), uuid.New())

	assert.NoError(t, err)
}

func TestActiveClients_MultipleClientsSameChannel(t *testing.T) {
	channelID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	mockRepo := &mockParticipantRepository{
		conversationIDs: map[uuid.UUID][]uuid.UUID{
			userID1: {channelID},
			userID2: {channelID},
		},
	}
	ac := NewActiveClients(context.Background(), mockRepo)

	client1 := &Client{
		Id:     uuid.New(),
		UserID: userID1,
	}
	client2 := &Client{
		Id:     uuid.New(),
		UserID: userID2,
	}

	ac.AddClient(client1)
	ac.AddClient(client2)

	clients := getClientsByChannel(ac, channelID)
	assert.Len(t, clients, 2)
}

func newBenchmarkActiveClients(conversationIDs map[uuid.UUID][]uuid.UUID) *activeClients {
	mockRepo := &mockParticipantRepository{
		conversationIDs: conversationIDs,
	}
	return NewActiveClients(context.Background(), mockRepo)
}

func newBenchmarkClient(userID uuid.UUID) *Client {
	return &Client{
		Id:          uuid.New(),
		UserID:      userID,
		sendChannel: make(chan OutgoingNotification, 16),
	}
}

func BenchmarkActiveClients_AddClient(b *testing.B) {
	ac := newBenchmarkActiveClients(make(map[uuid.UUID][]uuid.UUID))
	clients := make([]*Client, b.N)
	for i := 0; i < b.N; i++ {
		clients[i] = newBenchmarkClient(uuid.New())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ac.AddClient(clients[i])
	}
}

func BenchmarkActiveClients_RemoveClient(b *testing.B) {
	ac := newBenchmarkActiveClients(make(map[uuid.UUID][]uuid.UUID))
	clients := make([]*Client, b.N)
	for i := 0; i < b.N; i++ {
		client := newBenchmarkClient(uuid.New())
		clients[i] = client
		ac.AddClient(client)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ac.RemoveClient(clients[i])
	}
}

func BenchmarkActiveClients_InvalidateMembership(b *testing.B) {
	userID := uuid.New()
	conversationID := uuid.New()
	ac := newBenchmarkActiveClients(map[uuid.UUID][]uuid.UUID{
		userID: {conversationID},
	})
	client := newBenchmarkClient(userID)
	ac.AddClient(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ac.InvalidateMembership(context.Background(), userID)
	}
}

func BenchmarkActiveClients_NotifyChannelClients(b *testing.B) {
	ac := newBenchmarkActiveClients(make(map[uuid.UUID][]uuid.UUID))
	channelID := uuid.New()
	for i := 0; i < 50; i++ {
		userID := uuid.New()
		client := newBenchmarkClient(userID)
		ac.participants = &mockParticipantRepository{
			conversationIDs: map[uuid.UUID][]uuid.UUID{
				userID: {channelID},
			},
		}
		ac.AddClient(client)
	}
	notification := OutgoingNotification{
		Type:   "benchmark",
		UserID: uuid.New(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ac.NotifyChannelClients(context.Background(), channelID, notification)
		clients := getClientsByChannel(ac, channelID)
		for _, client := range clients {
			for len(client.sendChannel) > 0 {
				<-client.sendChannel
			}
		}
	}
}
