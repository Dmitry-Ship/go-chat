package services

import (
	"GitHub/go-chat/backend/pkg/domain"

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
	conversations        domain.ConversationCommandRepository
	participants         domain.ParticipantCommandRepository
	messagingService     MessagingService
	notificationsService NotificationsService
}

func NewConversationService(
	conversations domain.ConversationCommandRepository,
	participants domain.ParticipantCommandRepository,
	messagingService MessagingService,
	notificationsService NotificationsService,
) *conversationService {
	return &conversationService{
		conversations:        conversations,
		participants:         participants,
		messagingService:     messagingService,
		notificationsService: notificationsService,
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

	err = s.participants.Store(domain.NewParticipant(newConversationID, fromUserId))

	if err != nil {
		return uuid.Nil, err
	}

	err = s.participants.Store(domain.NewParticipant(newConversationID, toUserId))

	if err != nil {
		return uuid.Nil, err
	}

	err = s.notificationsService.SubscribeToTopic("conversation:"+newConversationID.String(), fromUserId)

	if err != nil {
		return uuid.Nil, err
	}

	err = s.notificationsService.SubscribeToTopic("conversation:"+newConversationID.String(), toUserId)

	if err != nil {
		return uuid.Nil, err
	}

	return newConversationID, nil
}

func (s *conversationService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation := domain.NewPublicConversation(id, name)
	err := s.conversations.StorePublicConversation(conversation)

	if err != nil {
		return err
	}

	err = s.JoinPublicConversation(conversation.ID, userId)

	return err
}

func (s *conversationService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.Store(domain.NewParticipant(conversationID, userId))

	if err != nil {
		return err
	}

	err = s.notificationsService.SubscribeToTopic("conversation:"+conversationID.String(), userId)

	if err != nil {
		return err
	}

	err = s.messagingService.SendJoinedConversationMessage(conversationID, userId)

	return err
}

func (s *conversationService) RenamePublicConversation(conversationID uuid.UUID, userId uuid.UUID, name string) error {
	conversation, err := s.conversations.GetPublicConversation(conversationID)

	if err != nil {
		return err
	}

	conversation.Rename(name)

	err = s.conversations.UpdatePublicConversation(conversation)

	if err != nil {
		return err
	}

	go s.notificationsService.NotifyAboutConversationRenamed(conversationID, name)

	err = s.messagingService.SendRenamedConversationMessage(conversationID, userId, name)

	return err
}

func (s *conversationService) LeavePublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	err := s.participants.DeleteByConversationIDAndUserID(conversationID, userId)

	if err != nil {
		return err
	}

	err = s.notificationsService.UnsubscribeFromTopic("conversation:"+conversationID.String(), userId)

	if err != nil {
		return err
	}

	err = s.messagingService.SendLeftConversationMessage(conversationID, userId)

	return err
}

func (s *conversationService) DeleteConversation(id uuid.UUID) error {
	go s.notificationsService.NotifyAboutConversationDeletion(id)

	err := s.conversations.Delete(id)

	return err
}
