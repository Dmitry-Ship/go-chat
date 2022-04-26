package services

import (
	"GitHub/go-chat/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

type messagesRepositoryMock struct {
	methodCalled string
}

func (m *messagesRepositoryMock) StoreTextMessage(message *domain.TextMessage) error {
	m.methodCalled = "StoreTextMessage"
	return nil
}

func (m *messagesRepositoryMock) StoreJoinedConversationMessage(message *domain.JoinedConversationMessage) error {
	m.methodCalled = "StoreJoinedConversationMessage"
	return nil
}

func (m *messagesRepositoryMock) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	m.methodCalled = "StoreRenamedConversationMessage"
	return nil
}

func (m *messagesRepositoryMock) StoreLeftConversationMessage(message *domain.LeftConversationMessage) error {
	m.methodCalled = "StoreLeftConversationMessage"
	return nil
}

type pubSubMock struct {
	methodCalled string
}

func (m *pubSubMock) Subscribe(topic string) <-chan domain.DomainEvent {
	m.methodCalled = "Subscribe"
	return nil
}

func (m *pubSubMock) Publish(event domain.DomainEvent) {
	m.methodCalled = "Publish"
}

func (m *pubSubMock) Close() {
	m.methodCalled = "Close"
}

func TestNewMessagingService(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{}
	pubSub := &pubSubMock{}

	messagingService := NewMessagingService(messagesRepository, pubSub)

	if messagingService == nil {
		t.Error("Expected messagingService to be not nil")
	}

	if messagingService.messages != messagesRepository {
		t.Error("Expected messagingService.messagesRepository to be equal to messagesRepository")
	}

}

func TestMessagingService_SendTextMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{}
	pubSub := &pubSubMock{}
	messagingService := NewMessagingService(messagesRepository, pubSub)

	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendTextMessage("test", conversationID, userID)

	if err != nil {
		t.Error("Expected err to be nil")
	}

	if messagesRepository.methodCalled != "StoreTextMessage" {
		t.Error("Expected StoreTextMessage to be called")
	}
}

func TestMessagingService_SendJoinedConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{}
	pubSub := &pubSubMock{}
	messagingService := NewMessagingService(messagesRepository, pubSub)

	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendJoinedConversationMessage(conversationID, userID)

	if err != nil {
		t.Error("Expected err to be nil")
	}

	if messagesRepository.methodCalled != "StoreJoinedConversationMessage" {
		t.Error("Expected StoreJoinedConversationMessage to be called")
	}

	if pubSub.methodCalled != "Publish" {
		t.Error("Expected Publish to be called")
	}
}

func TestMessagingService_SendRenamedConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{}
	pubSub := &pubSubMock{}
	messagingService := NewMessagingService(messagesRepository, pubSub)

	conversationID := uuid.New()
	newName := "new name"
	userID := uuid.New()

	err := messagingService.SendRenamedConversationMessage(conversationID, userID, newName)

	if err != nil {
		t.Error("Expected err to be nil")
	}

	if messagesRepository.methodCalled != "StoreRenamedConversationMessage" {
		t.Error("Expected StoreRenamedConversationMessage to be called")
	}

	if pubSub.methodCalled != "Publish" {
		t.Error("Expected Publish to be called")
	}
}

func TestMessagingService_SendLeftConversationMessage(t *testing.T) {
	messagesRepository := &messagesRepositoryMock{}
	pubSub := &pubSubMock{}
	messagingService := NewMessagingService(messagesRepository, pubSub)
	conversationID := uuid.New()
	userID := uuid.New()

	err := messagingService.SendLeftConversationMessage(conversationID, userID)

	if err != nil {
		t.Error("Expected err to be nil")
	}

	if messagesRepository.methodCalled != "StoreLeftConversationMessage" {
		t.Error("Expected StoreLeftConversationMessage to be called")
	}

	if pubSub.methodCalled != "Publish" {
		t.Error("Expected Publish to be called")
	}
}
