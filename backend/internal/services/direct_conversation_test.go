package services

import (
	"context"
	"testing"

	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDirectConversationRepository struct {
	mock.Mock
}

func (m *MockDirectConversationRepository) Store(ctx context.Context, conversation *domain.DirectConversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockDirectConversationRepository) GetID(ctx context.Context, firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, firstUserID, secondUserID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockDirectConversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DirectConversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DirectConversation), args.Error(1)
}

type MockNotificationServiceForDirect struct {
	mock.Mock
}

func (m *MockNotificationServiceForDirect) Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	args := m.Called(ctx, conversationID, notification)
	return args.Error(0)
}

func (m *MockNotificationServiceForDirect) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID {
	args := m.Called(ctx, conn, userID)
	return args.Get(0).(uuid.UUID)
}

func (m *MockNotificationServiceForDirect) Run() {
	m.Called()
}

func (m *MockNotificationServiceForDirect) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationServiceForDirect) Shutdown() {
	m.Called()
}

type MockCacheServiceForDirect struct {
	mock.Mock
}

func (m *MockCacheServiceForDirect) InvalidateConversation(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *MockCacheServiceForDirect) InvalidateParticipants(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *MockCacheServiceForDirect) InvalidateUserConversations(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestDirectConversationService_StartDirectConversation(t *testing.T) {
	ctx := context.Background()
	fromUserID := uuid.New()
	toUserID := uuid.New()
	existingConversationID := uuid.New()

	mockDirectConversations := new(MockDirectConversationRepository)
	mockNotifications := new(MockNotificationServiceForDirect)
	mockCache := new(MockCacheServiceForDirect)

	service := NewDirectConversationService(
		mockDirectConversations,
		mockNotifications,
		mockCache,
	)

	t.Run("existing conversation returns existing ID", func(t *testing.T) {
		mockDirectConversations.ExpectedCalls = nil
		mockDirectConversations.On("GetID", mock.Anything, fromUserID, toUserID).Return(existingConversationID, nil)

		conversationID, err := service.StartDirectConversation(ctx, fromUserID, toUserID)

		assert.NoError(t, err)
		assert.Equal(t, existingConversationID, conversationID)
		mockDirectConversations.AssertExpectations(t)
	})

	t.Run("successful new conversation", func(t *testing.T) {
		mockDirectConversations.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockDirectConversations.On("GetID", mock.Anything, fromUserID, toUserID).Return(uuid.Nil, assert.AnError)
		mockDirectConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.DirectConversation")).Return(nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, fromUserID).Return(nil)
		mockCache.On("InvalidateUserConversations", mock.Anything, fromUserID).Return(nil)

		conversationID, err := service.StartDirectConversation(ctx, fromUserID, toUserID)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, conversationID)
		mockDirectConversations.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockDirectConversations.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockDirectConversations.On("GetID", mock.Anything, fromUserID, toUserID).Return(uuid.Nil, assert.AnError)
		mockDirectConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.DirectConversation")).Return(assert.AnError)

		conversationID, err := service.StartDirectConversation(ctx, fromUserID, toUserID)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, conversationID)
	})

	t.Run("notification error", func(t *testing.T) {
		mockDirectConversations.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockDirectConversations.On("GetID", mock.Anything, fromUserID, toUserID).Return(uuid.Nil, assert.AnError)
		mockDirectConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.DirectConversation")).Return(nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, fromUserID).Return(assert.AnError)

		conversationID, err := service.StartDirectConversation(ctx, fromUserID, toUserID)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, conversationID)
	})

	t.Run("cache error", func(t *testing.T) {
		mockDirectConversations.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockDirectConversations.On("GetID", mock.Anything, fromUserID, toUserID).Return(uuid.Nil, assert.AnError)
		mockDirectConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.DirectConversation")).Return(nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, fromUserID).Return(nil)
		mockCache.On("InvalidateUserConversations", mock.Anything, fromUserID).Return(assert.AnError)

		conversationID, err := service.StartDirectConversation(ctx, fromUserID, toUserID)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, conversationID)
	})
}
