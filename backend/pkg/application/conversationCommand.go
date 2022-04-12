package application

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"

	"github.com/google/uuid"
)

type ConversationCommandService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	DeleteConversation(id uuid.UUID) error
	SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	RenamePublicConversation(conversationId uuid.UUID, userId uuid.UUID, name string) error
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

	message := domain.NewJoinedConversationMessage(conversationID, userId)

	err = s.messages.StoreJoinedConversation(message)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) RenamePublicConversation(conversationID uuid.UUID, userId uuid.UUID, name string) error {
	err := s.conversations.RenameConversation(conversationID, name)

	if err != nil {
		return err
	}

	message := domain.NewConversationRenamedMessage(conversationID, userId, name)

	err = s.messages.StoreRenamedConversation(message)

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

	message := domain.NewLeftConversationMessage(conversationID, userId)

	err = s.messages.StoreLeftConversation(message)

	if err != nil {
		return err
	}

	s.conversationWSResolver.DispatchSystemMessage(message.ID, conversationID)

	return nil
}

func (s *conversationCommandService) DeleteConversation(id uuid.UUID) error {
	s.conversationWSResolver.DispatchConversationDeleted(id)

	err := s.conversations.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *conversationCommandService) SendTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewTextMessage(conversationId, userId, messageText)

	err := s.messages.StoreTextMessage(message)

	if err != nil {
		return err
	}

	s.conversationWSResolver.DispatchUserMessage(message.ID, conversationId, userId)

	return nil
}
