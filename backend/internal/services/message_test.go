package services

import (
	"context"
	"testing"

	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQueriesRepositoryForMessage struct {
	mock.Mock
}

func (m *MockQueriesRepositoryForMessage) GetContacts(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(userID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(conversationID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(conversationID, userID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) GetUserByID(userID uuid.UUID) (readModel.UserDTO, error) {
	args := m.Called(userID)
	return args.Get(0).(readModel.UserDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) GetConversation(id uuid.UUID, userID uuid.UUID) (readModel.ConversationFullDTO, error) {
	args := m.Called(id, userID)
	return args.Get(0).(readModel.ConversationFullDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ConversationDTO, error) {
	args := m.Called(userID, paginationInfo)
	return args.Get(0).([]readModel.ConversationDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) RenameConversationAndReturn(conversationID uuid.UUID, name string) error {
	args := m.Called(conversationID, name)
	return args.Error(0)
}

func (m *MockQueriesRepositoryForMessage) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.MessageDTO, error) {
	args := m.Called(conversationID, requestUserID, paginationInfo)
	return args.Get(0).([]readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	args := m.Called(messageID, requestUserID)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) StoreMessageAndReturnWithUser(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (readModel.MessageDTO, error) {
	args := m.Called(id, conversationID, userID, content, messageType)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) StoreSystemMessageAndReturn(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (readModel.MessageDTO, error) {
	args := m.Called(id, conversationID, userID, content, messageType)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) IsMember(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) IsMemberOwner(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) InviteToConversationAtomic(conversationID uuid.UUID, inviteeID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(conversationID, inviteeID, participantID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) KickParticipantAtomic(conversationID uuid.UUID, targetID uuid.UUID) (int64, error) {
	args := m.Called(conversationID, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueriesRepositoryForMessage) LeaveConversationAtomic(conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
	args := m.Called(conversationID, userID)
	return args.Get(0).(int64), args.Error(1)
}

type MockNotificationServiceForMessage struct {
	mock.Mock
}

func (m *MockNotificationServiceForMessage) Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	args := m.Called(ctx, conversationID, notification)
	return args.Error(0)
}

func (m *MockNotificationServiceForMessage) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID {
	args := m.Called(ctx, conn, userID)
	return args.Get(0).(uuid.UUID)
}

func (m *MockNotificationServiceForMessage) Run() {
	m.Called()
}

func (m *MockNotificationServiceForMessage) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationServiceForMessage) Shutdown() {
	m.Called()
}

func TestMessageService_SendTextMessage(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()
	messageText := "Hello, world!"

	mockQueries := new(MockQueriesRepositoryForMessage)
	mockNotifications := new(MockNotificationServiceForMessage)

	service := NewMessageService(
		mockQueries,
		mockNotifications,
	)

	t.Run("successful send", func(t *testing.T) {
		mockQueries.On("IsMember", conversationID, userID).Return(true, nil)
		mockQueries.On("StoreMessageAndReturnWithUser", mock.Anything, conversationID, userID, messageText, int32(0)).Return(readModel.MessageDTO{}, nil)
		mockNotifications.On("Broadcast", mock.Anything, conversationID, mock.Anything).Return(nil)

		err := service.SendTextMessage(ctx, conversationID, userID, messageText)

		assert.NoError(t, err)
		mockQueries.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
	})

	t.Run("user not in conversation", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMember", conversationID, userID).Return(false, nil)

		err := service.SendTextMessage(ctx, conversationID, userID, messageText)

		assert.Error(t, err)
	})

	t.Run("store error", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		mockQueries.On("IsMember", conversationID, userID).Return(true, nil)
		mockQueries.On("StoreMessageAndReturnWithUser", mock.Anything, conversationID, userID, messageText, int32(0)).Return(readModel.MessageDTO{}, assert.AnError)

		err := service.SendTextMessage(ctx, conversationID, userID, messageText)

		assert.Error(t, err)
	})
}
