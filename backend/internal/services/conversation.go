package services

import (
	"fmt"

	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateGroupConversation(conversationID uuid.UUID, name string, userID uuid.UUID) error
	DeleteGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error
	Rename(conversationID uuid.UUID, userID uuid.UUID, name string) error
	Join(conversationID uuid.UUID, userID uuid.UUID) error
	Leave(conversationID uuid.UUID, userID uuid.UUID) error
	Invite(conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
	Kick(conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error
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
	users               domain.UserRepository
	messages            domain.MessageRepository
}

func NewConversationService(
	groupConversations domain.GroupConversationRepository,
	directConversations domain.DirectConversationRepository,
	participants domain.ParticipantRepository,
	users domain.UserRepository,
	messages domain.MessageRepository,
) *conversationService {
	return &conversationService{
		groupConversations:  groupConversations,
		directConversations: directConversations,
		participants:        participants,
		users:               users,
		messages:            messages,
	}
}

func (s *conversationService) CreateGroupConversation(conversationID uuid.UUID, name string, userID uuid.UUID) error {
	if err := domain.ValidateConversationName(name); err != nil {
		return fmt.Errorf("validate conversation name error: %w", err)
	}

	user, err := s.users.GetByID(userID)

	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	conversation, err := domain.NewGroupConversation(conversationID, name, *user)

	if err != nil {
		return fmt.Errorf("new group conversation error: %w", err)
	}

	if err := s.groupConversations.Store(conversation); err != nil {
		return fmt.Errorf("store conversation error: %w", err)
	}

	return nil
}

func (s *conversationService) StartDirectConversation(fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error) {
	existingConversationID, err := s.directConversations.GetID(fromUserID, toUserID)

	if err == nil {
		return existingConversationID, nil
	}

	newConversationID := uuid.New()

	conversation, err := domain.NewDirectConversation(newConversationID, toUserID, fromUserID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("new direct conversation error: %w", err)
	}

	if err = s.directConversations.Store(conversation); err != nil {
		return uuid.Nil, fmt.Errorf("store conversation error: %w", err)
	}

	return newConversationID, nil
}

func (s *conversationService) DeleteGroupConversation(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	if err = conversation.Delete(participant); err != nil {
		return fmt.Errorf("delete conversation error: %w", err)
	}

	if err := s.groupConversations.Update(conversation); err != nil {
		return fmt.Errorf("update conversation error: %w", err)
	}

	return nil
}

func (s *conversationService) Rename(conversationID uuid.UUID, userID uuid.UUID, name string) error {
	if err := domain.ValidateConversationName(name); err != nil {
		return fmt.Errorf("validate conversation name error: %w", err)
	}

	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	if err = conversation.Rename(name, participant); err != nil {
		return fmt.Errorf("rename conversation error: %w", err)
	}

	if err := s.groupConversations.Update(conversation); err != nil {
		return fmt.Errorf("update conversation error: %w", err)
	}

	return nil
}

func (s *conversationService) SendDirectTextMessage(conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	conversation, err := s.directConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	message, err := conversation.SendTextMessage(messageID, messageText, *participant)

	if err != nil {
		return fmt.Errorf("send text message error: %w", err)
	}

	if err := s.messages.Store(message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *conversationService) SendGroupTextMessage(conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	messageID := uuid.New()

	message, err := conversation.SendTextMessage(messageID, messageText, participant)

	if err != nil {
		return fmt.Errorf("send text message error: %w", err)
	}

	if err := s.messages.Store(message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *conversationService) SendJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	user, err := s.users.GetByID(userID)

	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	message, err := conversation.SendJoinedConversationMessage(messageID, user)

	if err != nil {
		return fmt.Errorf("send joined message error: %w", err)
	}

	if err := s.messages.Store(message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *conversationService) SendInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	user, err := s.users.GetByID(userID)

	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	message, err := conversation.SendInvitedConversationMessage(messageID, user)

	if err != nil {
		return fmt.Errorf("send invited message error: %w", err)
	}

	if err := s.messages.Store(message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *conversationService) SendRenamedConversationMessage(conversationID uuid.UUID, userID uuid.UUID, name string) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	message, err := conversation.SendRenamedConversationMessage(messageID, participant, name)

	if err != nil {
		return fmt.Errorf("send renamed message error: %w", err)
	}

	if err := s.messages.Store(message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *conversationService) SendLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	message, err := conversation.SendLeftConversationMessage(messageID, participant)

	if err != nil {
		return fmt.Errorf("send left message error: %w", err)
	}

	if err := s.messages.Store(message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *conversationService) Join(conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	user, err := s.users.GetByID(userID)

	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	participant, err := conversation.Join(*user)

	if err != nil {
		return fmt.Errorf("join conversation error: %w", err)
	}

	if err := s.participants.Store(participant); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	return nil
}

func (s *conversationService) Leave(conversationID uuid.UUID, userID uuid.UUID) error {
	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err = conversation.Leave(participant)

	if err != nil {
		return fmt.Errorf("leave conversation error: %w", err)
	}

	if err := s.participants.Update(participant); err != nil {
		return fmt.Errorf("update participant error: %w", err)
	}

	return nil
}

func (s *conversationService) Invite(conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	invitee, err := s.users.GetByID(inviteeID)

	if err != nil {
		return fmt.Errorf("get invitee by id error: %w", err)
	}

	newParticipant, err := conversation.Invite(participant, invitee)

	if err != nil {
		return fmt.Errorf("invite user error: %w", err)
	}

	if err := s.participants.Store(newParticipant); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	return nil
}

func (s *conversationService) Kick(conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	kicker, err := s.participants.GetByConversationIDAndUserID(conversationID, kickerID)

	if err != nil {
		return fmt.Errorf("get kicker participant error: %w", err)
	}

	target, err := s.participants.GetByConversationIDAndUserID(conversationID, targetID)

	if err != nil {
		return fmt.Errorf("get target participant error: %w", err)
	}

	kicked, err := conversation.Kick(kicker, target)

	if err != nil {
		return fmt.Errorf("kick user error: %w", err)
	}

	if err := s.participants.Update(kicked); err != nil {
		return fmt.Errorf("update participant error: %w", err)
	}

	return nil
}
