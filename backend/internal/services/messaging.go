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
	pubsub   domain.PubSub
}

func NewMessagingService(
	messages domain.MessageCommandRepository,
	pubsub domain.PubSub,
) *messagingService {
	return &messagingService{
		messages: messages,
		pubsub:   pubsub,
	}
}

func (s *messagingService) SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewTextMessage(conversationId, userId, messageText)

	err := s.messages.StoreTextMessage(message)

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewMessageSent(conversationId, message.ID, userId))

	return nil
}

func (s *messagingService) SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewJoinedConversationMessage(conversationId, userId)

	err := s.messages.StoreJoinedConversationMessage(message)

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewMessageSent(conversationId, message.ID, userId))

	return nil
}

func (s *messagingService) SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error {
	message := domain.NewConversationRenamedMessage(conversationId, userId, name)

	err := s.messages.StoreRenamedConversationMessage(message)

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewMessageSent(conversationId, message.ID, userId))

	return nil
}

func (s *messagingService) SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewLeftConversationMessage(conversationId, userId)

	err := s.messages.StoreLeftConversationMessage(message)

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewMessageSent(conversationId, message.ID, userId))

	return nil
}
