package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockCacheClient struct {
	deletedKeys     []string
	deletedPatterns []string
	deleteError     error
	setError        error
}

func (m *mockCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (m *mockCacheClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return m.setError
}

func (m *mockCacheClient) Delete(ctx context.Context, key string) error {
	m.deletedKeys = append(m.deletedKeys, key)
	return m.deleteError
}

func (m *mockCacheClient) DeletePattern(ctx context.Context, pattern string) error {
	m.deletedPatterns = append(m.deletedPatterns, pattern)
	return m.deleteError
}

func TestNewCacheService(t *testing.T) {
	mockClient := &mockCacheClient{}
	service := NewCacheService(mockClient)

	assert.NotNil(t, service)
	assert.IsType(t, &cacheService{}, service)
}

func TestCacheService_InvalidateConversation(t *testing.T) {
	mockClient := &mockCacheClient{}
	service := NewCacheService(mockClient)
	ctx := context.Background()
	conversationID := uuid.New()

	err := service.InvalidateConversation(ctx, conversationID)

	assert.NoError(t, err)
	assert.Len(t, mockClient.deletedKeys, 2)
	assert.Len(t, mockClient.deletedPatterns, 1)
	assert.Contains(t, mockClient.deletedKeys[0], conversationID.String())
	assert.Contains(t, mockClient.deletedKeys[1], conversationID.String())
	assert.Contains(t, mockClient.deletedPatterns[0], conversationID.String())
}

func TestCacheService_InvalidateConversation_Error(t *testing.T) {
	mockClient := &mockCacheClient{
		deleteError: assert.AnError,
	}
	service := NewCacheService(mockClient)
	ctx := context.Background()
	conversationID := uuid.New()

	err := service.InvalidateConversation(ctx, conversationID)

	assert.Error(t, err)
}

func TestCacheService_InvalidateParticipants(t *testing.T) {
	mockClient := &mockCacheClient{}
	service := NewCacheService(mockClient)
	ctx := context.Background()
	conversationID := uuid.New()

	err := service.InvalidateParticipants(ctx, conversationID)

	assert.NoError(t, err)
	assert.Len(t, mockClient.deletedKeys, 2)
	assert.Contains(t, mockClient.deletedKeys[0], conversationID.String())
	assert.Contains(t, mockClient.deletedKeys[1], conversationID.String())
}

func TestCacheService_InvalidateParticipants_Error(t *testing.T) {
	mockClient := &mockCacheClient{
		deleteError: assert.AnError,
	}
	service := NewCacheService(mockClient)
	ctx := context.Background()
	conversationID := uuid.New()

	err := service.InvalidateParticipants(ctx, conversationID)

	assert.Error(t, err)
}

func TestCacheService_InvalidateUserConversations(t *testing.T) {
	mockClient := &mockCacheClient{}
	service := NewCacheService(mockClient)
	ctx := context.Background()
	userID := uuid.New()

	err := service.InvalidateUserConversations(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, mockClient.deletedKeys, 1)
	assert.Contains(t, mockClient.deletedKeys[0], userID.String())
}

func TestCacheService_InvalidateUserConversations_Error(t *testing.T) {
	mockClient := &mockCacheClient{
		deleteError: assert.AnError,
	}
	service := NewCacheService(mockClient)
	ctx := context.Background()
	userID := uuid.New()

	err := service.InvalidateUserConversations(ctx, userID)

	assert.Error(t, err)
}
