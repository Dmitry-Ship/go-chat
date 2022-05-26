package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateGroupConversation(id uuid.UUID, name string, userID uuid.UUID) error
	DeleteGroupConversation(id uuid.UUID, userID uuid.UUID) error
	RenameGroupConversation(conversationID uuid.UUID, userID uuid.UUID, name string) error
	JoinGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error
	LeaveGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error
	InviteToGroupConversation(conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
	StartDirectConversation(fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error)
	SendDirectTextMessage(messageText string, conversationID uuid.UUID, userID uuid.UUID) error
	SendGroupTextMessage(messageText string, conversationID uuid.UUID, userID uuid.UUID) error
	SendJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error
	SendInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error
	SendRenamedConversationMessage(conversationID uuid.UUID, userID uuid.UUID, name string) error
	SendLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error
}

type conversationService struct {
	groupConversations  domain.GroupConversationRepository
	directConversations domain.DirectConversationRepository
	participants        domain.ParticipantRepository
	messages            domain.MessageRepository
}

func NewConversationService(
	groupConversations domain.GroupConversationRepository,
	directConversations domain.DirectConversationRepository,
	participants domain.ParticipantRepository,
	messages domain.MessageRepository,
) *conversationService {
	return &conversationService{
		groupConversations:  groupConversations,
		directConversations: directConversations,
		participants:        participants,
		messages:            messages,
	}
}

func (s *conversationService) CreateGroupConversation(id uuid.UUID, name string, userID uuid.UUID) error {
	conversation, err := domain.NewGroupConversation(id, name, userID)

	if err != nil {
		return err
	}

	return s.groupConversations.Store(conversation)
}

func (s *conversationService) StartDirectConversation(fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error) {
	existingConversationID, err := s.directConversations.GetID(fromUserID, toUserID)

	if err == nil {
		return existingConversationID, nil
	}

	newConversationID := uuid.New()

	conversation, err := domain.NewDirectConversation(newConversationID, toUserID, fromUserID)

	if err != nil {
		return uuid.Nil, err
	}

	err = s.directConversations.Store(conversation)

	if err != nil {
		return uuid.Nil, err
	}

	return newConversationID, nil
}

func (s *conversationService) DeleteGroupConversation(id uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(id)

	if err != nil {
		return err
	}

	err = conversation.Delete(userID)

	if err != nil {
		return err
	}

	return s.groupConversations.Update(conversation)
}

func (s *conversationService) RenameGroupConversation(conversationID uuid.UUID, userID uuid.UUID, name string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	err = conversation.Rename(name, userID)

	if err != nil {
		return err
	}

	return s.groupConversations.Update(conversation)
}

func (s *conversationService) SendDirectTextMessage(messageText string, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.directConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	message, err := conversation.SendTextMessage(messageText, userID)

	if err != nil {
		return err
	}

	return s.messages.StoreTextMessage(message)
}

func (s *conversationService) SendGroupTextMessage(messageText string, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return err
	}

	message, err := conversation.SendTextMessage(messageText, participant)

	if err != nil {
		return err
	}

	return s.messages.StoreTextMessage(message)
}

func (s *conversationService) SendJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	message, err := conversation.SendJoinedConversationMessage(conversationID, userID)

	if err != nil {
		return err
	}

	return s.messages.StoreJoinedConversationMessage(message)
}

func (s *conversationService) SendInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	message, err := conversation.SendInvitedConversationMessage(conversationID, userID)

	if err != nil {
		return err
	}

	return s.messages.StoreInvitedConversationMessage(message)
}

func (s *conversationService) SendRenamedConversationMessage(conversationID uuid.UUID, userID uuid.UUID, name string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	message, err := conversation.SendRenamedConversationMessage(conversationID, userID, name)

	if err != nil {
		return err
	}

	return s.messages.StoreRenamedConversationMessage(message)
}

func (s *conversationService) SendLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	message, err := conversation.SendLeftConversationMessage(conversationID, userID)

	if err != nil {
		return err
	}

	return s.messages.StoreLeftConversationMessage(message)
}

func (s *conversationService) JoinGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	participant, err := conversation.Join(userID)

	if err != nil {
		return err
	}

	return s.participants.Store(participant)
}

func (s *conversationService) LeaveGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error {
	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return err
	}

	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	participant, err = conversation.Leave(participant)

	if err != nil {
		return err
	}

	return s.participants.Update(participant)
}

func (s *conversationService) InviteToGroupConversation(conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	participant, err := conversation.Invite(userID, inviteeID)

	if err != nil {
		return err
	}

	return s.participants.Store(participant)
}
