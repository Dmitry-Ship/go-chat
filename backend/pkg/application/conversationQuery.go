package application

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
)

type conversationDataDTO struct {
	Conversation domain.ConversationDTO `json:"conversation"`
	Joined       bool                   `json:"joined"`
}

type MessageFullDTO struct {
	*domain.MessageDTO
	User      *domain.UserDTO `json:"user,omitempty"`
	IsInbound bool            `json:"is_inbound,omitempty"`
}

type ConversationQueryService interface {
	GetConversation(conversationId uuid.UUID, userId uuid.UUID) (conversationDataDTO, error)
	GetConversations() ([]*domain.ConversationDTO, error)
	GetConversationMessages(conversationId uuid.UUID, userId uuid.UUID) ([]MessageFullDTO, error)
}

type conversationQueryService struct {
	conversations domain.ConversationQueryRepository
	participants  domain.ParticipantCommandRepository
	users         domain.UserCommandRepository
	messages      domain.MessageQueryRepository
}

func NewConversationQueryService(
	conversations domain.ConversationQueryRepository,
	participants domain.ParticipantCommandRepository,
	users domain.UserCommandRepository,
	messages domain.MessageQueryRepository,
) *conversationQueryService {
	return &conversationQueryService{
		conversations: conversations,
		users:         users,
		participants:  participants,
		messages:      messages,
	}
}

func (s *conversationQueryService) GetConversation(conversationId uuid.UUID, userId uuid.UUID) (conversationDataDTO, error) {
	conversation, err := s.conversations.FindByID(conversationId)

	if err != nil {
		return conversationDataDTO{}, err
	}

	data := conversationDataDTO{
		Conversation: *conversation,
		Joined:       s.hasJoined(conversationId, userId),
	}

	return data, nil
}

func (s *conversationQueryService) GetConversations() ([]*domain.ConversationDTO, error) {
	return s.conversations.FindAll()

}

func (s *conversationQueryService) hasJoined(conversationID uuid.UUID, userId uuid.UUID) bool {
	_, err := s.participants.FindByConversationIDAndUserID(conversationID, userId)

	return err == nil
}

func (s *conversationQueryService) GetConversationMessages(conversationId uuid.UUID, userID uuid.UUID) ([]MessageFullDTO, error) {
	messages, err := s.messages.FindAllByConversationID(conversationId)

	if err != nil {
		return nil, err
	}

	messagesFull := []MessageFullDTO{}

	for _, message := range messages {

		var messageFull MessageFullDTO

		messageFull.MessageDTO = message

		if message.Type == "user" {
			user, err := s.users.FindByID(*message.UserId)

			if err != nil {
				return nil, err
			}

			messageFull.User = user
			messageFull.IsInbound = *message.UserId != userID
		}

		messagesFull = append(messagesFull, messageFull)
	}

	return messagesFull, nil
}
