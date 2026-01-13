package services

import (
	"context"
	"testing"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGroupConversationRepository struct {
	mock.Mock
}

func (m *MockGroupConversationRepository) Store(ctx context.Context, conversation *domain.GroupConversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockGroupConversationRepository) Update(ctx context.Context, conversation *domain.GroupConversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockGroupConversationRepository) Rename(ctx context.Context, id uuid.UUID, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockGroupConversationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGroupConversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GroupConversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GroupConversation), args.Error(1)
}

type MockQueriesRepository struct {
	mock.Mock
}

func (m *MockQueriesRepository) GetContacts(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(userID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(conversationID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(conversationID, userID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetUserByID(userID uuid.UUID) (readModel.UserDTO, error) {
	args := m.Called(userID)
	return args.Get(0).(readModel.UserDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetUsersByIDs(userIDs []uuid.UUID) ([]readModel.UserDTO, error) {
	args := m.Called(userIDs)
	return args.Get(0).([]readModel.UserDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetConversation(id uuid.UUID, userID uuid.UUID) (readModel.ConversationFullDTO, error) {
	args := m.Called(id, userID)
	return args.Get(0).(readModel.ConversationFullDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ConversationDTO, error) {
	args := m.Called(userID, paginationInfo)
	return args.Get(0).([]readModel.ConversationDTO), args.Error(1)
}

func (m *MockQueriesRepository) RenameConversationAndReturn(conversationID uuid.UUID, name string) error {
	args := m.Called(conversationID, name)
	return args.Error(0)
}

func (m *MockQueriesRepository) GetConversationMessages(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.MessageDTO, error) {
	args := m.Called(conversationID, paginationInfo)
	return args.Get(0).([]readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepository) GetNotificationMessage(messageID uuid.UUID) (readModel.MessageDTO, error) {
	args := m.Called(messageID)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepository) StoreMessageAndReturn(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (readModel.MessageDTO, error) {
	args := m.Called(id, conversationID, userID, content, messageType)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepository) IsMember(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesRepository) IsMemberOwner(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesRepository) InviteToConversationAtomic(conversationID uuid.UUID, inviteeID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(conversationID, inviteeID, participantID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockQueriesRepository) LeaveConversationAtomic(conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
	args := m.Called(conversationID, userID)
	return args.Get(0).(int64), args.Error(1)
}

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error) {
	args := m.Called(ctx, message)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	args := m.Called(ctx, conversationID, notification)
	return args.Error(0)
}

func (m *MockNotificationService) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID {
	args := m.Called(ctx, conn, userID)
	return args.Get(0).(uuid.UUID)
}

func (m *MockNotificationService) Run() {
	m.Called()
}

func (m *MockNotificationService) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationService) Shutdown() {
	m.Called()
}

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) InvalidateConversation(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateParticipants(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateUserConversations(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestGroupConversationService_CreateGroupConversation(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()
	name := "Test Group"

	mockGroupConversations := new(MockGroupConversationRepository)
	mockQueries := new(MockQueriesRepository)
	mockMessages := new(MockMessageService)
	mockNotifications := new(MockNotificationService)
	mockCache := new(MockCacheService)

	service := NewGroupConversationService(
		mockGroupConversations,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful creation", func(t *testing.T) {
		mockGroupConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.GroupConversation")).Return(nil)
		mockCache.On("InvalidateUserConversations", mock.Anything, userID).Return(nil)

		err := service.CreateGroupConversation(ctx, conversationID, name, userID)

		assert.NoError(t, err)
		mockGroupConversations.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockGroupConversations.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockGroupConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.GroupConversation")).Return(assert.AnError)

		err := service.CreateGroupConversation(ctx, conversationID, name, userID)

		assert.Error(t, err)
	})

	t.Run("cache error", func(t *testing.T) {
		mockGroupConversations.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockGroupConversations.On("Store", mock.Anything, mock.AnythingOfType("*domain.GroupConversation")).Return(nil)
		mockCache.On("InvalidateUserConversations", mock.Anything, userID).Return(assert.AnError)

		err := service.CreateGroupConversation(ctx, conversationID, name, userID)

		assert.Error(t, err)
	})
}

func TestGroupConversationService_DeleteGroupConversation(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()

	mockGroupConversations := new(MockGroupConversationRepository)
	mockQueries := new(MockQueriesRepository)
	mockMessages := new(MockMessageService)
	mockNotifications := new(MockNotificationService)
	mockCache := new(MockCacheService)

	service := NewGroupConversationService(
		mockGroupConversations,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful deletion", func(t *testing.T) {
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(true, nil)
		mockGroupConversations.On("Delete", mock.Anything, conversationID).Return(nil)
		mockCache.On("InvalidateConversation", mock.Anything, conversationID).Return(nil)
		mockNotifications.On("Broadcast", mock.Anything, conversationID, mock.Anything).Return(nil)

		err := service.DeleteGroupConversation(ctx, conversationID, userID)

		assert.NoError(t, err)
		mockQueries.AssertExpectations(t)
		mockGroupConversations.AssertExpectations(t)
		mockCache.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
	})

	t.Run("user not owner", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(false, nil)

		err := service.DeleteGroupConversation(ctx, conversationID, userID)

		assert.Error(t, err)
	})
}

func TestGroupConversationService_Rename(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()
	newName := "New Name"

	mockGroupConversations := new(MockGroupConversationRepository)
	mockQueries := new(MockQueriesRepository)
	mockMessages := new(MockMessageService)
	mockNotifications := new(MockNotificationService)
	mockCache := new(MockCacheService)

	service := NewGroupConversationService(
		mockGroupConversations,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful rename", func(t *testing.T) {
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(true, nil)
		mockGroupConversations.On("Rename", mock.Anything, conversationID, newName).Return(nil)
		mockMessages.On("Send", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(readModel.MessageDTO{}, nil)
		mockCache.On("InvalidateConversation", mock.Anything, conversationID).Return(nil)
		mockQueries.On("GetConversation", conversationID, userID).Return(readModel.ConversationFullDTO{ID: conversationID}, nil)
		mockNotifications.On("Broadcast", mock.Anything, conversationID, mock.Anything).Return(nil)

		err := service.Rename(ctx, conversationID, userID, newName)

		assert.NoError(t, err)
		mockQueries.AssertExpectations(t)
		mockGroupConversations.AssertExpectations(t)
		mockMessages.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("user not owner", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(false, nil)

		err := service.Rename(ctx, conversationID, userID, newName)

		assert.Error(t, err)
	})
}
