package services

import (
	"GitHub/go-chat/backend/pkg/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	DeleteConversation(id uuid.UUID) error
	SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	RenamePublicConversation(conversationId uuid.UUID, userId uuid.UUID, name string) error
}

type conversationService struct {
	conversations          domain.ConversationCommandRepository
	participants           domain.ParticipantCommandRepository
	messages               domain.MessageCommandRepository
	conversationWSResolver ConversationWSResolver
}

func NewConversationService(
	conversations domain.ConversationCommandRepository,
	participants domain.ParticipantCommandRepository,
	messages domain.MessageCommandRepository,
	conversationWSResolver ConversationWSResolver,
) *conversationService {
	return &conversationService{
		conversations:          conversations,
		participants:           participants,
		messages:               messages,
		conversationWSResolver: conversationWSResolver,
	}
}

func (s *conversationService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewPublicConversation(id, name)
	err := s.conversations.StorePublicConversation(conversation)

	if err != nil {
		return err
	}

	err = s.JoinPublicConversation(conversation.ID, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(conversationID, userId))

	if err != nil {
		return err
	}

	message := domain.NewJoinedConversationMessage(conversationID, userId)

	err = s.messages.StoreJoinedConversationMessage(message)

	if err != nil {
		return err
	}

	go s.conversationWSResolver.NotifyAboutMessage(conversationID, message.ID, userId)

	return nil
}

func (s *conversationService) RenamePublicConversation(conversationID uuid.UUID, userId uuid.UUID, name string) error {
	err := s.conversations.RenamePublicConversation(conversationID, name)

	if err != nil {
		return err
	}

	message := domain.NewConversationRenamedMessage(conversationID, userId, name)

	err = s.messages.StoreRenamedConversationMessage(message)

	if err != nil {
		return err
	}

	go s.conversationWSResolver.NotifyAboutMessage(conversationID, message.ID, userId)
	go s.conversationWSResolver.NotifyAboutConversationRenamed(conversationID, name)

	return nil
}

func (s *conversationService) LeavePublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	message := domain.NewLeftConversationMessage(conversationID, userId)

	err = s.messages.StoreLeftConversationMessage(message)

	if err != nil {
		return err
	}

	go s.conversationWSResolver.NotifyAboutMessage(conversationID, message.ID, userId)

	return nil
}

func (s *conversationService) DeleteConversation(id uuid.UUID) error {
	go s.conversationWSResolver.NotifyAboutConversationDeletion(id)

	err := s.conversations.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationService) SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewTextMessage(conversationId, userId, messageText)

	err := s.messages.StoreTextMessage(message)

	if err != nil {
		return err
	}

	go s.conversationWSResolver.NotifyAboutMessage(conversationId, message.ID, userId)

	return nil
}
