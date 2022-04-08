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
	conversations domain.ConversationRepository
	participants  domain.ParticipantRepository
	users         domain.UserRepository
	messages      domain.MessageRepository
}

func NewConversationQueryService(conversations domain.ConversationRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.MessageRepository) *conversationQueryService {
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
		messageFull, err := s.makeMessageFullDTO(message, userID)

		if err != nil {
			return nil, err
		}

		messagesFull = append(messagesFull, messageFull)
	}

	return messagesFull, nil
}

func (s *conversationQueryService) makeMessageFullDTO(message *domain.MessageDTO, userID uuid.UUID) (MessageFullDTO, error) {
	if message.Type == "system" {
		return MessageFullDTO{
			MessageDTO: message,
		}, nil
	}

	user, err := s.users.FindByID(userID)

	if err != nil {
		return MessageFullDTO{}, err
	}

	m := MessageFullDTO{
		MessageDTO: message,
		User:       user,
		IsInbound:  user.ID != userID,
	}

	return m, nil
}
