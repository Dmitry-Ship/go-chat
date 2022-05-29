package services

import (
	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateGroupConversation(conversationID uuid.UUID, name string, userID uuid.UUID) error
	DeleteGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error
	RenameGroupConversation(conversationID uuid.UUID, userID uuid.UUID, name string) error
	JoinGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error
	LeaveGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error
	InviteToGroupConversation(conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
	StartDirectConversation(fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error)
	SendDirectTextMessage(conversationID uuid.UUID, userID uuid.UUID, messageText string) error
	SendGroupTextMessage(conversationID uuid.UUID, userID uuid.UUID, messageText string) error
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

func (s *conversationService) CreateGroupConversation(conversationID uuid.UUID, name string, userID uuid.UUID) error {
	conversationName, err := domain.NewConversationName(name)

	if err != nil {
		return err
	}

	conversation, err := domain.NewGroupConversation(conversationID, conversationName, userID)

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

	if err = s.directConversations.Store(conversation); err != nil {
		return uuid.Nil, err
	}

	return newConversationID, nil
}

func (s *conversationService) DeleteGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	if err = conversation.Delete(userID); err != nil {
		return err
	}

	return s.groupConversations.Update(conversation)
}

func (s *conversationService) RenameGroupConversation(conversationID uuid.UUID, userID uuid.UUID, name string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	conversationName, err := domain.NewConversationName(name)

	if err != nil {
		return err
	}

	if err = conversation.Rename(conversationName, userID); err != nil {
		return err
	}

	return s.groupConversations.Update(conversation)
}

func (s *conversationService) SendDirectTextMessage(conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	conversation, err := s.directConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	messageID := uuid.New()

	message, err := conversation.SendTextMessage(messageID, messageText, userID)

	if err != nil {
		return err
	}

	return s.messages.Store(message)
}

func (s *conversationService) SendGroupTextMessage(conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return err
	}

	messageID := uuid.New()

	message, err := conversation.SendTextMessage(messageID, messageText, participant)

	if err != nil {
		return err
	}

	return s.messages.Store(message)
}

func (s *conversationService) SendJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	messageID := uuid.New()

	message, err := conversation.SendJoinedConversationMessage(messageID, userID)

	if err != nil {
		return err
	}

	return s.messages.Store(message)
}

func (s *conversationService) SendInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	messageID := uuid.New()

	message, err := conversation.SendInvitedConversationMessage(messageID, userID)

	if err != nil {
		return err
	}

	return s.messages.Store(message)
}

func (s *conversationService) SendRenamedConversationMessage(conversationID uuid.UUID, userID uuid.UUID, name string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	messageID := uuid.New()

	message, err := conversation.SendRenamedConversationMessage(messageID, userID, name)

	if err != nil {
		return err
	}

	return s.messages.Store(message)
}

func (s *conversationService) SendLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return err
	}

	messageID := uuid.New()

	message, err := conversation.SendLeftConversationMessage(messageID, userID)

	if err != nil {
		return err
	}

	return s.messages.Store(message)
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
