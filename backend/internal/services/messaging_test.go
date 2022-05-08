package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type messagesRepositoryMock struct {
	methodsCalled map[string]int
}

func (m *messagesRepositoryMock) StoreTextMessage(message *domain.TextMessage) error {
	m.methodsCalled["StoreTextMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreJoinedConversationMessage(message *domain.JoinedConversationMessage) error {
	m.methodsCalled["StoreJoinedConversationMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreInvitedConversationMessage(message *domain.JoinedConversationMessage) error {
	m.methodsCalled["StoreInvitedConversationMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	m.methodsCalled["StoreRenamedConversationMessage"]++
	return nil
}

func (m *messagesRepositoryMock) StoreLeftConversationMessage(message *domain.LeftConversationMessage) error {
	m.methodsCalled["StoreLeftConversationMessage"]++
	return nil
}

func TestNewMessagingService(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}

	messagingService := NewMessagingService(messagesRepository)

	assert.Equal(t, messagesRepository, messagingService.messages)
}

func TestMessagingService_SendTextMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagingService := NewMessagingService(messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendTextMessage("test", conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreTextMessage"])
}

func TestMessagingService_SendJoinedConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagingService := NewMessagingService(messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendJoinedConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreJoinedConversationMessage"])
}

func TestMessagingService_SendInvitedConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagingService := NewMessagingService(messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendInvitedConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreInvitedConversationMessage"])
}

func TestMessagingService_SendRenamedConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagingService := NewMessagingService(messagesRepository)
	conversationID := uuid.New()
	newName := "new name"
	userID := uuid.New()

	err := messagingService.SendRenamedConversationMessage(conversationID, userID, newName)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreRenamedConversationMessage"])
}

func TestMessagingService_SendLeftConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{
		methodsCalled: make(map[string]int),
	}
	messagingService := NewMessagingService(messagesRepository)
	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendLeftConversationMessage(conversationID, userID)

	assert.Nil(t, err)
	assert.Equal(t, 1, messagesRepository.methodsCalled["StoreLeftConversationMessage"])
}
