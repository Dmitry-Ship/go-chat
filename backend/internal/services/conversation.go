package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	StartPrivateConversation(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error)
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	DeletePublicConversation(id uuid.UUID, userId uuid.UUID) error
	RenamePublicConversation(conversationId uuid.UUID, userId uuid.UUID, name string) error
}

type conversationService struct {
	conversations domain.ConversationRepository
	participants  domain.ParticipantRepository
}

func NewConversationService(
	conversations domain.ConversationRepository,
	participants domain.ParticipantRepository,
) *conversationService {
	return &conversationService{
		conversations: conversations,
		participants:  participants,
	}
}

func (s *conversationService) StartPrivateConversation(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error) {
	existingConversationID, err := s.conversations.GetPrivateConversationID(fromUserId, toUserId)

	if err == nil {
		return existingConversationID, nil
	}

	newConversationID := uuid.New()

	conversation := domain.NewPrivateConversation(newConversationID, toUserId, fromUserId)

	err = s.conversations.StorePrivateConversation(conversation)

	if err != nil {
		return uuid.Nil, err
	}

	return newConversationID, nil
}

func (s *conversationService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewPublicConversation(id, name, userId)

	return s.conversations.StorePublicConversation(conversation)
}

func (s *conversationService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	return s.participants.Store(domain.NewJoinedParticipant(conversationID, userId))
}

func (s *conversationService) RenamePublicConversation(conversationID uuid.UUID, userId uuid.UUID, name string) error {
	conversation, err := s.conversations.GetPublicConversation(conversationID)

	if err != nil {
		return err
	}

	err = conversation.Rename(name, userId)

	if err != nil {
		return err
	}

	return s.conversations.UpdatePublicConversation(conversation)
}

func (s *conversationService) LeavePublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	err = participant.LeavePublicConversation(conversationID)

	if err != nil {
		return err
	}

	return s.participants.Update(participant)
}

func (s *conversationService) DeletePublicConversation(id uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.conversations.GetPublicConversation(id)

	if err != nil {
		return err
	}

	err = conversation.Delete(userID)

	if err != nil {
		return err
	}

	return s.conversations.UpdatePublicConversation(conversation)
}
