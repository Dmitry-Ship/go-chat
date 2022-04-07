package application

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
)

type conversationData struct {
	Conversation domain.Conversation `json:"conversation"`
	Joined       bool                `json:"joined"`
}

type MessageFull struct {
	User *domain.User `json:"user"`
	*domain.Message
	IsInbound bool `json:"is_inbound"`
}

type ConversationQueryService interface {
	GetConversation(conversationId uuid.UUID, userId uuid.UUID) (conversationData, error)
	GetConversations() ([]*domain.Conversation, error)
	GetConversationMessages(conversationId uuid.UUID, userId uuid.UUID) ([]MessageFull, error)
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

func (s *conversationQueryService) GetConversation(conversationId uuid.UUID, userId uuid.UUID) (conversationData, error) {

	conversation, err := s.conversations.FindByID(conversationId)

	if err != nil {
		return conversationData{}, err
	}

	data := conversationData{
		Conversation: *conversation,
		Joined:       s.hasJoined(conversationId, userId),
	}

	return data, nil
}

func (s *conversationQueryService) GetConversations() ([]*domain.Conversation, error) {
	return s.conversations.FindAll()
}

func (s *conversationQueryService) hasJoined(conversationID uuid.UUID, userId uuid.UUID) bool {
	_, err := s.participants.FindByConversationIDAndUserID(conversationID, userId)

	return err == nil
}

func (s *conversationQueryService) GetConversationMessages(conversationId uuid.UUID, userID uuid.UUID) ([]MessageFull, error) {
	messages, err := s.messages.FindAllByConversationID(conversationId)

	if err != nil {
		return nil, err
	}

	messagesFull := []MessageFull{}

	for _, message := range messages {
		messageFull, err := s.makeMessageFull(message, userID)

		if err != nil {
			return nil, err
		}

		messagesFull = append(messagesFull, messageFull)
	}

	return messagesFull, nil
}

func (s *conversationQueryService) makeMessageFull(message *domain.Message, userID uuid.UUID) (MessageFull, error) {
	user, err := s.users.FindByID(message.UserID)

	if err != nil {
		return MessageFull{}, err
	}

	m := MessageFull{
		User:      user,
		Message:   message,
		IsInbound: user.ID != userID,
	}

	return m, nil
}
