package services

import (
	"GitHub/go-chat/backend/pkg/domain"

	"github.com/google/uuid"
)

type MessagingService interface {
	SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
	SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error
	SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
}

type messagingService struct {
	messages             domain.MessageCommandRepository
	notificationsService NotificationsService
}

func NewMessagingService(
	messages domain.MessageCommandRepository,
	notificationsService NotificationsService,
) *messagingService {
	return &messagingService{
		messages:             messages,
		notificationsService: notificationsService,
	}
}

func (s *messagingService) SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewTextMessage(conversationId, userId, messageText)

	err := s.messages.StoreTextMessage(message)

	if err != nil {
		return err
	}

	go s.notificationsService.NotifyAboutMessage(conversationId, message.ID, userId)

	return nil
}

func (s *messagingService) SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewJoinedConversationMessage(conversationId, userId)

	err := s.messages.StoreJoinedConversationMessage(message)

	if err != nil {
		return err
	}

	go s.notificationsService.NotifyAboutMessage(conversationId, message.ID, userId)

	return nil
}

func (s *messagingService) SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error {
	message := domain.NewConversationRenamedMessage(conversationId, userId, name)

	err := s.messages.StoreRenamedConversationMessage(message)

	if err != nil {
		return err
	}

	go s.notificationsService.NotifyAboutMessage(conversationId, message.ID, userId)

	return nil
}

func (s *messagingService) SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewLeftConversationMessage(conversationId, userId)

	err := s.messages.StoreLeftConversationMessage(message)

	if err != nil {
		return err
	}

	go s.notificationsService.NotifyAboutMessage(conversationId, message.ID, userId)

	return nil
}
