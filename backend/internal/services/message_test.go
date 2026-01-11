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

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Send(ctx context.Context, message *domain.Message, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	args := m.Called(ctx, message, requestUserID)
	return args.Get(0).(readModel.MessageDTO), args.Error(1)
}

type MockNotificationServiceForMessageTest struct {
	mock.Mock
}

func (m *MockNotificationServiceForMessageTest) Broadcast(ctx context.Context, conversationID uuid.UUID, notification ws.OutgoingNotification) error {
	args := m.Called(ctx, conversationID, notification)
	return args.Error(0)
}

func (m *MockNotificationServiceForMessageTest) RegisterClient(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) uuid.UUID {
	args := m.Called(ctx, conn, userID)
	return args.Get(0).(uuid.UUID)
}

func (m *MockNotificationServiceForMessageTest) Run() {
	m.Called()
}

func (m *MockNotificationServiceForMessageTest) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationServiceForMessageTest) Shutdown() {
	m.Called()
}

func TestMessageService_Send(t *testing.T) {
	ctx := context.Background()
	conversationID := uuid.New()
	userID := uuid.New()
	content := "Hello, world!"

	mockRepository := new(MockMessageRepository)
	mockNotifications := new(MockNotificationServiceForMessageTest)

	service := NewMessageService(mockRepository, mockNotifications)

	t.Run("successful send", func(t *testing.T) {
		message, err := domain.NewMessage(conversationID, userID, domain.MessageTypeUser, content)
		assert.NoError(t, err)

		mockRepository.On("Send", mock.Anything, message, userID).Return(readModel.MessageDTO{}, nil)
		mockNotifications.On("Broadcast", mock.Anything, conversationID, mock.Anything).Return(nil)

		_, err = service.Send(ctx, message, userID)

		assert.NoError(t, err)
		mockRepository.AssertExpectations(t)
		mockNotifications.AssertExpectations(t)
	})

	t.Run("send error", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		message, err := domain.NewMessage(conversationID, userID, domain.MessageTypeUser, content)
		assert.NoError(t, err)

		mockRepository.On("Send", mock.Anything, message, userID).Return(readModel.MessageDTO{}, assert.AnError)

		_, err = service.Send(ctx, message, userID)

		assert.Error(t, err)
	})

	t.Run("broadcast error", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockNotifications.ExpectedCalls = nil
		message, err := domain.NewMessage(conversationID, userID, domain.MessageTypeUser, content)
		assert.NoError(t, err)

		mockRepository.On("Send", mock.Anything, message, userID).Return(readModel.MessageDTO{}, nil)
		mockNotifications.On("Broadcast", mock.Anything, conversationID, mock.Anything).Return(assert.AnError)

		_, err = service.Send(ctx, message, userID)

		assert.Error(t, err)
	})
}
