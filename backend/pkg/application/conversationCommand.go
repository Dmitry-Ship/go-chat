package application

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
	"fmt"

	"github.com/google/uuid"
)

type ConversationCommandService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	DeleteConversation(id uuid.UUID) error
	SendUserMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	sendSystemMessage(messageText string, conversationId uuid.UUID) error
	RenamePublicConversation(conversationId uuid.UUID, name string) error
}

type conversationCommandService struct {
	conversations          domain.ConversationCommandRepository
	participants           domain.ParticipantCommandRepository
	users                  domain.UserCommandRepository
	usersQuery             readModel.UserQueryRepository
	messages               domain.MessageCommandRepository
	conversationWSResolver ConversationWSResolver
}

func NewConversationCommandService(
	conversations domain.ConversationCommandRepository,
	participants domain.ParticipantCommandRepository,
	users domain.UserCommandRepository,
	usersQuery readModel.UserQueryRepository,
	messages domain.MessageCommandRepository,
	conversationWSResolver ConversationWSResolver,
) *conversationCommandService {
	return &conversationCommandService{
		conversations:          conversations,
		users:                  users,
		usersQuery:             usersQuery,
		participants:           participants,
		messages:               messages,
		conversationWSResolver: conversationWSResolver,
	}
}

func (s *conversationCommandService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewConversation(id, name, false)
	err := s.conversations.Store(conversation)

	if err != nil {
		return err
	}

	err = s.JoinPublicConversation(conversation.ID, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(conversationID, userId))

	if err != nil {
		return err
	}

	// Dispatch event to send system message

	user, err := s.usersQuery.GetUserByID(userId)

	if err != nil {
		return err
	}

	err = s.sendSystemMessage(fmt.Sprintf("%s joined", user.Name), conversationID)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) RenamePublicConversation(conversationID uuid.UUID, name string) error {
	err := s.conversations.RenameConversation(conversationID, name)

	if err != nil {
		return err
	}

	// Dispatch event to send system message and send update notification to all participants

	err = s.sendSystemMessage(fmt.Sprintf("chat has been renamed to %s", name), conversationID)

	if err != nil {
		return err
	}

	s.conversationWSResolver.DispatchConversationRenamed(conversationID, name)

	return nil
}

func (s *conversationCommandService) LeavePublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	// Dispatch event to send system message

	user, err := s.usersQuery.GetUserByID(userId)

	if err != nil {
		return err
	}

	err = s.sendSystemMessage(fmt.Sprintf("%s left", user.Name), conversationID)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) DeleteConversation(id uuid.UUID) error {

	// Dispatch event to send system message
	s.conversationWSResolver.DispatchConversationDeleted(id)

	err := s.conversations.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) SendUserMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewUserMessage(messageText, conversationId, userId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	// Dispatch event to send system message
	s.conversationWSResolver.DispatchUserMessage(message.ID, conversationId, userId)

	return nil
}

func (s *conversationCommandService) sendSystemMessage(messageText string, conversationId uuid.UUID) error {
	message := domain.NewSystemMessage(messageText, conversationId)

	err := s.messages.Store(message)

	if err != nil {
		return err
	}

	s.conversationWSResolver.DispatchSystemMessage(message.ID, conversationId)

	return nil
}
