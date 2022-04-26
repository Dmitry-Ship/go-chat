package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	CreatePrivateConversationIfNotExists(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error)
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	DeleteConversation(id uuid.UUID) error
	RenamePublicConversation(conversationId uuid.UUID, userId uuid.UUID, name string) error
}

type conversationService struct {
	conversations domain.ConversationCommandRepository
	participants  domain.ParticipantCommandRepository
	pubsub        domain.PubSub
}

func NewConversationService(
	conversations domain.ConversationCommandRepository,
	participants domain.ParticipantCommandRepository,
	pubsub domain.PubSub,
) *conversationService {
	return &conversationService{
		conversations: conversations,
		participants:  participants,
		pubsub:        pubsub,
	}
}

func (s *conversationService) CreatePrivateConversationIfNotExists(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error) {
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

	s.pubsub.Publish(domain.NewPrivateConversationCreated(newConversationID, fromUserId, toUserId))

	return newConversationID, nil
}

func (s *conversationService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewPublicConversation(id, name, userId)

	s.pubsub.Publish(domain.NewPublicConversationCreated(id, userId))

	return s.conversations.StorePublicConversation(conversation)
}

func (s *conversationService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(conversationID, userId))

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewPublicConversationJoined(conversationID, userId))

	return nil
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

	err = s.conversations.UpdatePublicConversation(conversation)

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewPublicConversationRenamed(conversationID, userId, name))

	return nil
}

func (s *conversationService) LeavePublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	s.pubsub.Publish(domain.NewPublicConversationLeft(conversationID, userId))

	return nil
}

func (s *conversationService) DeleteConversation(id uuid.UUID) error {
	s.pubsub.Publish(domain.NewPublicConversationDeleted(id))

	return s.conversations.Delete(id)
}
