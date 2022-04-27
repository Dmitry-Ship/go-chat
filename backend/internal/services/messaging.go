package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type MessagingService interface {
	SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
	SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error
	SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
}

type messagingService struct {
	messages domain.MessageCommandRepository
}

func NewMessagingService(
	messages domain.MessageCommandRepository,
) *messagingService {
	return &messagingService{
		messages: messages,
	}
}

func (s *messagingService) SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewTextMessage(conversationId, userId, messageText)

	return s.messages.StoreTextMessage(message)
}

func (s *messagingService) SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewJoinedConversationMessage(conversationId, userId)

	return s.messages.StoreJoinedConversationMessage(message)

}

func (s *messagingService) SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error {
	message := domain.NewConversationRenamedMessage(conversationId, userId, name)

	return s.messages.StoreRenamedConversationMessage(message)
}

func (s *messagingService) SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewLeftConversationMessage(conversationId, userId)

	return s.messages.StoreLeftConversationMessage(message)
}
