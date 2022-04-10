package application

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
)

type ConversationQueryService interface {
	GetConversation(conversationId uuid.UUID, userId uuid.UUID) (*domain.ConversationDTOFull, error)
	GetConversations() ([]*domain.ConversationDTO, error)
	GetConversationMessages(conversationId uuid.UUID, userId uuid.UUID) ([]*domain.MessageDTO, error)
}

type conversationQueryService struct {
	conversations domain.ConversationQueryRepository
	messages      domain.MessageQueryRepository
}

func NewConversationQueryService(conversations domain.ConversationQueryRepository, messages domain.MessageQueryRepository) *conversationQueryService {
	return &conversationQueryService{
		conversations: conversations,
		messages:      messages,
	}
}

func (s *conversationQueryService) GetConversation(conversationId uuid.UUID, userId uuid.UUID) (*domain.ConversationDTOFull, error) {
	return s.conversations.GetConversationByID(conversationId, userId)
}

func (s *conversationQueryService) GetConversations() ([]*domain.ConversationDTO, error) {
	return s.conversations.FindAll()
}

func (s *conversationQueryService) GetConversationMessages(conversationId uuid.UUID, userID uuid.UUID) ([]*domain.MessageDTO, error) {
	return s.messages.FindAllByConversationID(conversationId, userID)
}
