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

type MockParticipantRepository struct {
	mock.Mock
}

func (m *MockParticipantRepository) Store(ctx context.Context, participant *domain.Participant) error {
	args := m.Called(ctx, participant)
	return args.Error(0)
}

func (m *MockParticipantRepository) Delete(ctx context.Context, participantID uuid.UUID) error {
	args := m.Called(ctx, participantID)
	return args.Error(0)
}

func (m *MockParticipantRepository) GetByConversationIDAndUserID(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	args := m.Called(ctx, conversationID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Participant), args.Error(1)
}

func (m *MockParticipantRepository) GetIDsByConversationID(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, conversationID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func (m *MockParticipantRepository) GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

type MockQueriesRepositoryForMembership struct {
	mock.Mock
}

func (m *MockQueriesRepositoryForMembership) GetContacts(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(userID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(conversationID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	args := m.Called(conversationID, userID, paginationInfo)
	return args.Get(0).([]readModel.ContactDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetUserByID(userID uuid.UUID) (readModel.UserDTO, error) {
	args := m.Called(userID)
	return args.Get(0).(readModel.UserDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetUsersByIDs(userIDs []uuid.UUID) ([]readModel.UserDTO, error) {
	args := m.Called(userIDs)
	return args.Get(0).([]readModel.UserDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetConversation(id uuid.UUID, userID uuid.UUID) (readModel.ConversationFullDTO, error) {
	args := m.Called(id, userID)
	return args.Get(0).(readModel.ConversationFullDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ConversationDTO, error) {
	args := m.Called(userID, paginationInfo)
	return args.Get(0).([]readModel.ConversationDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) RenameConversationAndReturn(conversationID uuid.UUID, name string) error {
	args := m.Called(conversationID, name)
	return args.Error(0)
}

func (m *MockQueriesRepositoryForMembership) GetConversationMessages(conversationID uuid.UUID, cursor *readModel.MessageCursor, limit int) (readModel.MessagePageDTO, error) {
	args := m.Called(conversationID, cursor, limit)
	return args.Get(0).(readModel.MessagePageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) GetNotificationMessage(messageID uuid.UUID) (readModel.MessageDTO, error) {
	args := m.Called(messageID)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) StoreMessageAndReturn(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (readModel.MessageDTO, error) {
	args := m.Called(id, conversationID, userID, content, messageType)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) IsMember(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) IsMemberOwner(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) InviteToConversationAtomic(conversationID uuid.UUID, inviteeID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(conversationID, inviteeID, participantID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockQueriesRepositoryForMembership) LeaveConversationAtomic(conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
	args := m.Called(conversationID, userID)
	return args.Get(0).(int64), args.Error(1)
}

type MockMessageRepositoryForMembership struct {
	mock.Mock
}

func (m *MockMessageRepositoryForMembership) Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error) {
	args := m.Called(ctx, message)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

type MockMessageServiceForMembership struct {
	mock.Mock
}

func (m *MockMessageServiceForMembership) Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error) {
	args := m.Called(ctx, message)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

type MockNotificationServiceForMembership struct {
	mock.Mock
}

func (m *MockNotificationServiceForMembership) Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	args := m.Called(ctx, conversationID, notification)
	return args.Error(0)
}

func (m *MockNotificationServiceForMembership) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID {
	args := m.Called(ctx, conn, userID)
	return args.Get(0).(uuid.UUID)
}

func (m *MockNotificationServiceForMembership) Run() {
	m.Called()
}

func (m *MockNotificationServiceForMembership) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationServiceForMembership) Shutdown() {
	m.Called()
}

type MockCacheServiceForMembership struct {
	mock.Mock
}

func (m *MockCacheServiceForMembership) InvalidateConversation(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *MockCacheServiceForMembership) InvalidateParticipants(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *MockCacheServiceForMembership) InvalidateUserConversations(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestMembershipService_Join(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()

	mockParticipants := new(MockParticipantRepository)
	mockQueries := new(MockQueriesRepositoryForMembership)
	mockMessages := new(MockMessageServiceForMembership)
	mockNotifications := new(MockNotificationServiceForMembership)
	mockCache := new(MockCacheServiceForMembership)

	service := NewMembershipService(
		mockParticipants,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful join", func(t *testing.T) {
		mockParticipants.On("Store", mock.Anything, mock.AnythingOfType("*domain.Participant")).Return(nil)
		mockMessages.On("Send", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(readModel.MessageDTO{}, nil)
		mockCache.On("InvalidateParticipants", mock.Anything, conversationID).Return(nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, userID).Return(nil)

		err := service.Join(ctx, conversationID, userID)

		assert.NoError(t, err)
		mockParticipants.AssertExpectations(t)
		mockMessages.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockParticipants.ExpectedCalls = nil
		mockMessages.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockParticipants.On("Store", mock.Anything, mock.AnythingOfType("*domain.Participant")).Return(assert.AnError)

		err := service.Join(ctx, conversationID, userID)

		assert.Error(t, err)
	})
}

func TestMembershipService_Leave(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()

	mockParticipants := new(MockParticipantRepository)
	mockQueries := new(MockQueriesRepositoryForMembership)
	mockMessages := new(MockMessageServiceForMembership)
	mockNotifications := new(MockNotificationServiceForMembership)
	mockCache := new(MockCacheServiceForMembership)

	service := NewMembershipService(
		mockParticipants,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful leave", func(t *testing.T) {
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(false, nil)
		mockQueries.On("LeaveConversationAtomic", conversationID, userID).Return(int64(1), nil)
		mockMessages.On("Send", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(readModel.MessageDTO{}, nil)
		mockCache.On("InvalidateParticipants", mock.Anything, conversationID).Return(nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, userID).Return(nil)

		err := service.Leave(ctx, conversationID, userID)

		assert.NoError(t, err)
		mockQueries.AssertExpectations(t)
		mockMessages.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("owner cannot leave", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(true, nil)

		err := service.Leave(ctx, conversationID, userID)

		assert.Error(t, err)
	})

	t.Run("user not in conversation", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMemberOwner", conversationID, userID).Return(false, nil)
		mockQueries.On("LeaveConversationAtomic", conversationID, userID).Return(int64(0), nil)

		err := service.Leave(ctx, conversationID, userID)

		assert.Error(t, err)
	})
}

func TestMembershipService_Invite(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()
	inviteeID := uuid.New()

	mockParticipants := new(MockParticipantRepository)
	mockQueries := new(MockQueriesRepositoryForMembership)
	mockMessages := new(MockMessageServiceForMembership)
	mockNotifications := new(MockNotificationServiceForMembership)
	mockCache := new(MockCacheServiceForMembership)

	service := NewMembershipService(
		mockParticipants,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful invite", func(t *testing.T) {
		mockQueries.On("IsMember", conversationID, userID).Return(true, nil)
		mockQueries.On("InviteToConversationAtomic", conversationID, inviteeID, mock.Anything).Return(uuid.New(), nil)
		mockMessages.On("Send", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(readModel.MessageDTO{}, nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, inviteeID).Return(nil)
		mockCache.On("InvalidateParticipants", mock.Anything, conversationID).Return(nil)

		err := service.Invite(ctx, conversationID, userID, inviteeID)

		assert.NoError(t, err)
		mockQueries.AssertExpectations(t)
		mockMessages.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("user not in conversation", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMember", conversationID, userID).Return(false, nil)

		err := service.Invite(ctx, conversationID, userID, inviteeID)

		assert.Error(t, err)
	})
}

func TestMembershipService_Kick(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	kickerID := uuid.New()
	targetID := uuid.New()

	mockParticipants := new(MockParticipantRepository)
	mockQueries := new(MockQueriesRepositoryForMembership)
	mockMessages := new(MockMessageServiceForMembership)
	mockNotifications := new(MockNotificationServiceForMembership)
	mockCache := new(MockCacheServiceForMembership)

	service := NewMembershipService(
		mockParticipants,
		mockQueries,
		mockMessages,
		mockNotifications,
		mockCache,
	)

	t.Run("successful kick", func(t *testing.T) {
		mockQueries.On("IsMemberOwner", conversationID, kickerID).Return(true, nil)
		mockQueries.On("IsMember", conversationID, targetID).Return(true, nil)
		mockQueries.On("LeaveConversationAtomic", conversationID, targetID).Return(int64(1), nil)
		mockMessages.On("Send", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(readModel.MessageDTO{}, nil)
		mockCache.On("InvalidateParticipants", mock.Anything, conversationID).Return(nil)
		mockQueries.On("GetConversation", conversationID, kickerID).Return(readModel.ConversationFullDTO{ID: conversationID}, nil)
		mockNotifications.On("Broadcast", mock.Anything, conversationID, mock.Anything).Return(nil)
		mockNotifications.On("InvalidateMembership", mock.Anything, targetID).Return(nil)

		err := service.Kick(ctx, conversationID, kickerID, targetID)

		assert.NoError(t, err)
		mockQueries.AssertExpectations(t)
		mockMessages.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("kicker not owner", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMemberOwner", conversationID, kickerID).Return(false, nil)

		err := service.Kick(ctx, conversationID, kickerID, targetID)

		assert.Error(t, err)
	})

	t.Run("target not in conversation", func(t *testing.T) {
		mockQueries.ExpectedCalls = nil
		mockQueries.On("IsMemberOwner", conversationID, kickerID).Return(true, nil)
		mockQueries.On("IsMember", conversationID, targetID).Return(false, nil)

		err := service.Kick(ctx, conversationID, kickerID, targetID)

		assert.Error(t, err)
	})
}
