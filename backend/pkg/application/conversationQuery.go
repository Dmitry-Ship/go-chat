package application

import (
	"GitHub/go-chat/backend/pkg/readModel"

	"github.com/google/uuid"
)

type ConversationQueryService interface {
	GetConversation(conversationId uuid.UUID, userId uuid.UUID) (*readModel.ConversationFullDTO, error)
	GetConversations() ([]*readModel.ConversationDTO, error)
	GetConversationMessages(conversationId uuid.UUID, userId uuid.UUID) ([]*readModel.MessageDTO, error)
}

type conversationQueryService struct {
	conversations readModel.ConversationQueryRepository
	messages      readModel.MessageQueryRepository
}

func NewConversationQueryService(conversations readModel.ConversationQueryRepository, messages readModel.MessageQueryRepository) *conversationQueryService {
	return &conversationQueryService{
		conversations: conversations,
		messages:      messages,
	}
}

func (s *conversationQueryService) GetConversation(conversationId uuid.UUID, userId uuid.UUID) (*readModel.ConversationFullDTO, error) {
	return s.conversations.GetConversationByID(conversationId, userId)
}

func (s *conversationQueryService) GetConversations() ([]*readModel.ConversationDTO, error) {
	return s.conversations.FindAll()
}

func (s *conversationQueryService) GetConversationMessages(conversationId uuid.UUID, userID uuid.UUID) ([]*readModel.MessageDTO, error) {
	return s.messages.FindAllByConversationID(conversationId, userID)
}
