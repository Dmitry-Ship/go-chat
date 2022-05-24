package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error
	DeletePublicConversation(id uuid.UUID, userId uuid.UUID) error
	RenamePublicConversation(conversationId uuid.UUID, userId uuid.UUID, name string) error
	JoinPublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	LeavePublicConversation(conversationId uuid.UUID, userId uuid.UUID) error
	InviteToPublicConversation(conversationId uuid.UUID, userId uuid.UUID, inviteeID uuid.UUID) error
	StartPrivateConversation(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error)
	SendPrivateTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	SendPublicTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error
	SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
	SendInvitedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
	SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error
	SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error
}

type conversationService struct {
	publicConversations  domain.PublicConversationRepository
	privateConversations domain.PrivateConversationRepository
	participants         domain.ParticipantRepository
	messages             domain.MessageRepository
}

func NewConversationService(
	publicConversations domain.PublicConversationRepository,
	privateConversations domain.PrivateConversationRepository,
	participants domain.ParticipantRepository,
	messages domain.MessageRepository,
) *conversationService {
	return &conversationService{
		publicConversations:  publicConversations,
		privateConversations: privateConversations,
		participants:         participants,
		messages:             messages,
	}
}

func (s *conversationService) CreatePublicConversation(id uuid.UUID, name string, userId uuid.UUID) error {
	conversation, err := domain.NewPublicConversation(id, name, userId)

	if err != nil {
		return err
	}

	return s.publicConversations.Store(conversation)
}

func (s *conversationService) StartPrivateConversation(fromUserId uuid.UUID, toUserId uuid.UUID) (uuid.UUID, error) {
	existingConversationID, err := s.privateConversations.GetID(fromUserId, toUserId)

	if err == nil {
		return existingConversationID, nil
	}

	newConversationID := uuid.New()

	conversation, err := domain.NewPrivateConversation(newConversationID, toUserId, fromUserId)

	if err != nil {
		return uuid.Nil, err
	}

	err = s.privateConversations.Store(conversation)

	if err != nil {
		return uuid.Nil, err
	}

	return newConversationID, nil
}

func (s *conversationService) DeletePublicConversation(id uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.publicConversations.GetByID(id)

	if err != nil {
		return err
	}

	err = conversation.Delete(userID)

	if err != nil {
		return err
	}

	return s.publicConversations.Update(conversation)
}

func (s *conversationService) RenamePublicConversation(conversationID uuid.UUID, userId uuid.UUID, name string) error {
	conversation, err := s.publicConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	err = conversation.Rename(name, userId)

	if err != nil {
		return err
	}

	return s.publicConversations.Update(conversation)
}

func (s *conversationService) SendPrivateTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	conversation, err := s.privateConversations.GetByID(conversationId)

	if err != nil {
		return err
	}

	message, err := conversation.SendTextMessage(messageText, userId)

	if err != nil {
		return err
	}

	return s.messages.StoreTextMessage(message)
}

func (s *conversationService) SendPublicTextMessage(messageText string, conversationId uuid.UUID, userId uuid.UUID) error {
	conversation, err := s.publicConversations.GetByID(conversationId)

	if err != nil {
		return err
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationId, userId)

	if err != nil {
		return err
	}

	message, err := conversation.SendTextMessage(messageText, participant)

	if err != nil {
		return err
	}

	return s.messages.StoreTextMessage(message)
}

func (s *conversationService) SendJoinedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewJoinedConversationMessage(conversationId, userId)

	return s.messages.StoreJoinedConversationMessage(message)
}

func (s *conversationService) SendInvitedConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewInvitedConversationMessage(conversationId, userId)

	return s.messages.StoreInvitedConversationMessage(message)
}

func (s *conversationService) SendRenamedConversationMessage(conversationId uuid.UUID, userId uuid.UUID, name string) error {
	message := domain.NewConversationRenamedMessage(conversationId, userId, name)

	return s.messages.StoreRenamedConversationMessage(message)
}

func (s *conversationService) SendLeftConversationMessage(conversationId uuid.UUID, userId uuid.UUID) error {
	message := domain.NewLeftConversationMessage(conversationId, userId)

	return s.messages.StoreLeftConversationMessage(message)
}

func (s *conversationService) JoinPublicConversation(conversationID uuid.UUID, userId uuid.UUID) error {
	return s.participants.Store(domain.NewJoinedParticipant(conversationID, userId))
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

func (s *conversationService) InviteToPublicConversation(conversationID uuid.UUID, userId uuid.UUID, inviteeID uuid.UUID) error {
	newParticipant := domain.NewInvitedParticipant(conversationID, inviteeID)

	return s.participants.Store(newParticipant)
}
