package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type PublicConversationEditingService interface {
	Create(id uuid.UUID, name string, userId uuid.UUID) error
	Delete(id uuid.UUID, userId uuid.UUID) error
	Rename(conversationId uuid.UUID, userId uuid.UUID, name string) error
}

type publicConversationEditingService struct {
	conversations domain.PublicConversationRepository
}

func NewPublicConversationEditingService(
	conversations domain.PublicConversationRepository,
) *publicConversationEditingService {
	return &publicConversationEditingService{
		conversations: conversations,
	}
}

func (s *publicConversationEditingService) Create(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewPublicConversation(id, name, userId)

	return s.conversations.Store(conversation)
}

func (s *publicConversationEditingService) Delete(id uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.conversations.GetByID(id)

	if err != nil {
		return err
	}

	err = conversation.Delete(userID)

	if err != nil {
		return err
	}

	return s.conversations.Update(conversation)
}

func (s *publicConversationEditingService) Rename(conversationID uuid.UUID, userId uuid.UUID, name string) error {
	conversation, err := s.conversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	err = conversation.Rename(name, userId)

	if err != nil {
		return err
	}

	return s.conversations.Update(conversation)
}
