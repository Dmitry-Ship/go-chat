package application

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"

	"github.com/google/uuid"
)

type conversationDataDTO struct {
	Conversation mappers.ConversationDTO `json:"conversation"`
	Joined       bool                    `json:"joined"`
}

type MessageFullDTO struct {
	*mappers.MessageDTO
	User      *mappers.UserDTO `json:"user,omitempty"`
	IsInbound bool             `json:"is_inbound,omitempty"`
}

type ConversationQueryService interface {
	GetConversation(conversationId uuid.UUID, userId uuid.UUID) (conversationDataDTO, error)
	GetConversations() ([]*mappers.ConversationDTO, error)
	GetConversationMessages(conversationId uuid.UUID, userId uuid.UUID) ([]MessageFullDTO, error)
}

type conversationQueryService struct {
	conversations domain.ConversationRepository
	participants  domain.ParticipantRepository
	users         domain.UserRepository
	messages      domain.ChatMessageRepository
}

func NewConversationQueryService(conversations domain.ConversationRepository, participants domain.ParticipantRepository, users domain.UserRepository, messages domain.ChatMessageRepository) *conversationQueryService {
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
		Conversation: *mappers.ToConversationDTO(conversation),
		Joined:       s.hasJoined(conversationId, userId),
	}

	return data, nil
}

func (s *conversationQueryService) GetConversations() ([]*mappers.ConversationDTO, error) {

	conversations, err := s.conversations.FindAll()

	if err != nil {
		return nil, err
	}

	var result []*mappers.ConversationDTO

	for _, conversation := range conversations {
		result = append(result, mappers.ToConversationDTO(conversation))
	}

	return result, nil
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

func (s *conversationQueryService) makeMessageFullDTO(message *domain.Message, userID uuid.UUID) (MessageFullDTO, error) {
	if message.UserID == nil {
		return MessageFullDTO{
			MessageDTO: mappers.ToMessageDTO(message),
		}, nil
	}

	user, err := s.users.FindByID(*message.UserID)

	if err != nil {
		return MessageFullDTO{}, err
	}

	m := MessageFullDTO{
		MessageDTO: mappers.ToMessageDTO(message),
		User:       mappers.ToUserDTO(user),
		IsInbound:  user.ID != userID,
	}

	return m, nil
}
