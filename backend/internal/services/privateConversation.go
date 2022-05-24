package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type PrivateConversationService interface {
	Start(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error)
}

type privateConversationService struct {
	conversations domain.PrivateConversationRepository
}

func NewPrivateConversationService(
	conversations domain.PrivateConversationRepository,
) *privateConversationService {
	return &privateConversationService{
		conversations: conversations,
	}
}

func (s *privateConversationService) Start(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error) {
	existingConversationID, err := s.conversations.GetID(fromUserId, toUserId)

	if err == nil {
		return existingConversationID, nil
	}

	newConversationID := uuid.New()

	conversation, err := domain.NewPrivateConversation(newConversationID, toUserId, fromUserId)

	if err != nil {
		return uuid.Nil, err
	}

	err = s.conversations.Store(conversation)

	if err != nil {
		return uuid.Nil, err
	}

	return newConversationID, nil
}
