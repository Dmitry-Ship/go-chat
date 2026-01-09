package services

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateGroupConversation(ctx context.Context, conversationID uuid.UUID, name string, userID uuid.UUID) error
	DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Rename(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, name string) error
	Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
	Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error
	StartDirectConversation(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error)
	SendDirectTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error
	SendGroupTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error
}

type conversationService struct {
	groupConversations  domain.GroupConversationRepository
	directConversations domain.DirectConversationRepository
	participants        domain.ParticipantRepository
	users               domain.UserRepository
	messages            domain.MessageRepository
	notifications       NotificationService
	cache               CacheService
	systemMessages      SystemMessageService
	queries             readModel.QueriesRepository
}

func NewConversationService(
	groupConversations domain.GroupConversationRepository,
	directConversations domain.DirectConversationRepository,
	participants domain.ParticipantRepository,
	users domain.UserRepository,
	messages domain.MessageRepository,
	notifications NotificationService,
	cache CacheService,
	systemMessages SystemMessageService,
	queries readModel.QueriesRepository,
) *conversationService {
	return &conversationService{
		groupConversations:  groupConversations,
		directConversations: directConversations,
		participants:        participants,
		users:               users,
		messages:            messages,
		notifications:       notifications,
		cache:               cache,
		systemMessages:      systemMessages,
		queries:             queries,
	}
}

func (s *conversationService) CreateGroupConversation(ctx context.Context, conversationID uuid.UUID, name string, userID uuid.UUID) error {
	if err := domain.ValidateConversationName(name); err != nil {
		return fmt.Errorf("validate conversation name error: %w", err)
	}

	conversation, err := domain.NewGroupConversation(conversationID, name, userID)

	if err != nil {
		return fmt.Errorf("new group conversation error: %w", err)
	}

	if err := s.groupConversations.Store(ctx, conversation); err != nil {
		return fmt.Errorf("store conversation error: %w", err)
	}

	if err := s.cache.InvalidateUserConversations(ctx, userID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	return nil
}

func (s *conversationService) StartDirectConversation(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error) {
	existingConversationID, err := s.directConversations.GetID(ctx, fromUserID, toUserID)
	if err == nil {
		return existingConversationID, nil
	}

	newConversationID := uuid.New()

	conversation, err := domain.NewDirectConversation(newConversationID, toUserID, fromUserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("new direct conversation error: %w", err)
	}

	if err = s.directConversations.Store(ctx, conversation); err != nil {
		return uuid.Nil, fmt.Errorf("store conversation error: %w", err)
	}

	if err := s.cache.InvalidateUserConversations(ctx, fromUserID); err != nil {
		return uuid.Nil, fmt.Errorf("invalidate cache error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, fromUserID); err != nil {
		return uuid.Nil, fmt.Errorf("invalidate membership error: %w", err)
	}

	return newConversationID, nil
}

func (s *conversationService) DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)
	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	if err = conversation.Delete(participant); err != nil {
		return fmt.Errorf("delete conversation error: %w", err)
	}

	if err := s.groupConversations.Update(ctx, conversation); err != nil {
		return fmt.Errorf("update conversation error: %w", err)
	}

	if err := s.cache.InvalidateConversation(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "conversation_deleted", Payload: map[string]interface{}{"conversation_id": conversationID}}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) Rename(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, name string) error {
	if err := domain.ValidateConversationName(name); err != nil {
		return fmt.Errorf("validate conversation name error: %w", err)
	}

	conversation, err := s.groupConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	if err = conversation.Rename(name, participant); err != nil {
		return fmt.Errorf("rename conversation error: %w", err)
	}

	if err := s.groupConversations.Update(ctx, conversation); err != nil {
		return fmt.Errorf("update conversation error: %w", err)
	}

	if err := s.systemMessages.SaveRenamedMessage(ctx, conversationID, userID, name); err != nil {
		return fmt.Errorf("send renamed message error: %w", err)
	}

	if err := s.cache.InvalidateConversation(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	conversationDTO, err := s.queries.GetConversation(conversationID, userID)
	if err != nil {
		return fmt.Errorf("get conversation error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationDTO.ID, ws.OutgoingNotification{Type: "conversation_updated", Payload: conversationDTO}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) SendDirectTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	conversation, err := s.directConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	message, err := conversation.SendTextMessage(messageID, messageText, *participant)

	if err != nil {
		return fmt.Errorf("send text message error: %w", err)
	}

	if err := s.messages.Store(ctx, message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	messageDTO, err := s.queries.GetNotificationMessage(messageID, userID)
	if err != nil {
		return fmt.Errorf("get notification message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) SendGroupTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	messageID := uuid.New()

	message, err := conversation.SendTextMessage(messageID, messageText, participant)

	if err != nil {
		return fmt.Errorf("send text message error: %w", err)
	}

	if err := s.messages.Store(ctx, message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	messageDTO, err := s.queries.GetNotificationMessage(messageID, userID)
	if err != nil {
		return fmt.Errorf("get notification message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	user, err := s.users.GetByID(ctx, userID)

	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	participant, err := conversation.Join(*user)

	if err != nil {
		return fmt.Errorf("join conversation error: %w", err)
	}

	if err := s.participants.Store(ctx, participant); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	if err := s.systemMessages.SaveJoinedMessage(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("send joined message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	conversationDTO, err := s.queries.GetConversation(conversationID, userID)
	if err != nil {
		return fmt.Errorf("get conversation error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationDTO.ID, ws.OutgoingNotification{Type: "conversation_updated", Payload: conversationDTO}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, userID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *conversationService) Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	conversation, err := s.groupConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err = conversation.Leave(participant)

	if err != nil {
		return fmt.Errorf("leave conversation error: %w", err)
	}

	if err := s.participants.Delete(ctx, participant.ID); err != nil {
		return fmt.Errorf("delete participant error: %w", err)
	}

	if err := s.systemMessages.SaveLeftMessage(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("send left message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	conversationDTO, err := s.queries.GetConversation(conversationID, userID)
	if err != nil {
		return fmt.Errorf("get conversation error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationDTO.ID, ws.OutgoingNotification{Type: "conversation_updated", Payload: conversationDTO}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, userID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *conversationService) Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)

	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	invitee, err := s.users.GetByID(ctx, inviteeID)

	if err != nil {
		return fmt.Errorf("get invitee by id error: %w", err)
	}

	newParticipant, err := conversation.Invite(participant, invitee)

	if err != nil {
		return fmt.Errorf("invite user error: %w", err)
	}

	if err := s.participants.Store(ctx, newParticipant); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	if err := s.systemMessages.SaveInvitedMessage(ctx, conversationID, inviteeID); err != nil {
		return fmt.Errorf("send invited message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	conversationDTO, err := s.queries.GetConversation(conversationID, userID)
	if err != nil {
		return fmt.Errorf("get conversation error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationDTO.ID, ws.OutgoingNotification{Type: "conversation_updated", Payload: conversationDTO}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, inviteeID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *conversationService) Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)

	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	kicker, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, kickerID)

	if err != nil {
		return fmt.Errorf("get kicker participant error: %w", err)
	}

	target, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, targetID)

	if err != nil {
		return fmt.Errorf("get target participant error: %w", err)
	}

	kicked, err := conversation.Kick(kicker, target)

	if err != nil {
		return fmt.Errorf("kick user error: %w", err)
	}

	if err := s.participants.Delete(ctx, kicked.ID); err != nil {
		return fmt.Errorf("delete participant error: %w", err)
	}

	if err := s.systemMessages.SaveLeftMessage(ctx, conversationID, targetID); err != nil {
		return fmt.Errorf("send left message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	conversationDTO, err := s.queries.GetConversation(conversationID, kickerID)
	if err != nil {
		return fmt.Errorf("get conversation error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationDTO.ID, ws.OutgoingNotification{Type: "conversation_updated", Payload: conversationDTO}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, targetID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}
