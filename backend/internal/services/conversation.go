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

	if err := s.notifications.InvalidateMembership(ctx, fromUserID); err != nil {
		return uuid.Nil, fmt.Errorf("invalidate membership error: %w", err)
	}

	if err := s.cache.InvalidateUserConversations(ctx, fromUserID); err != nil {
		return uuid.Nil, fmt.Errorf("invalidate cache error: %w", err)
	}

	return newConversationID, nil
}

func (s *conversationService) DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	isOwner, err := s.queries.IsMemberOwner(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is owner error: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("user is not owner: %w", domain.ErrorUserNotOwner)
	}

	if err := s.groupConversations.Delete(ctx, conversationID); err != nil {
		return fmt.Errorf("delete conversation error: %w", err)
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

	isOwner, err := s.queries.IsMemberOwner(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is owner error: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("user is not owner: %w", domain.ErrorUserNotOwner)
	}

	if err := s.groupConversations.Rename(ctx, conversationID, name); err != nil {
		return fmt.Errorf("rename conversation error: %w", err)
	}

	if err := s.saveRenamedMessage(ctx, conversationID, userID, name); err != nil {
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
	isMember, err := s.queries.IsMember(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is member error: %w", err)
	}
	if !isMember {
		return fmt.Errorf("user is not in conversation: %w", domain.ErrorUserNotInConversation)
	}

	messageID := uuid.New()

	if _, err := domain.NewTextMessageContent(messageText); err != nil {
		return fmt.Errorf("validate message text error: %w", err)
	}

	messageDTO, err := s.queries.StoreMessageAndReturnWithUser(messageID, conversationID, userID, messageText, 0)
	if err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) SendGroupTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	isMember, err := s.queries.IsMember(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is member error: %w", err)
	}
	if !isMember {
		return fmt.Errorf("user is not in conversation: %w", domain.ErrorUserNotInConversation)
	}

	messageID := uuid.New()

	if _, err := domain.NewTextMessageContent(messageText); err != nil {
		return fmt.Errorf("validate message text error: %w", err)
	}

	messageDTO, err := s.queries.StoreMessageAndReturnWithUser(messageID, conversationID, userID, messageText, 0)
	if err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	participant := domain.NewParticipant(uuid.New(), conversationID, userID)

	if err := s.participants.Store(ctx, participant); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	if err := s.saveJoinedMessage(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("send joined message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, userID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *conversationService) Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	isOwner, err := s.queries.IsMemberOwner(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is owner error: %w", err)
	}
	if isOwner {
		return fmt.Errorf("owner cannot leave: %w", domain.ErrorOwnerCannotLeave)
	}

	rowsAffected, err := s.queries.LeaveConversationAtomic(conversationID, userID)
	if err != nil {
		return fmt.Errorf("leave conversation error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user is not in conversation: %w", domain.ErrorUserNotInConversation)
	}

	if err := s.saveLeftMessage(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("send left message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, userID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *conversationService) Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error {
	isMember, err := s.queries.IsMember(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is member error: %w", err)
	}
	if !isMember {
		return fmt.Errorf("user is not in conversation: %w", domain.ErrorUserNotInConversation)
	}

	participantID := uuid.New()

	storedInviteeID, err := s.queries.InviteToConversationAtomic(conversationID, inviteeID, participantID)
	if err != nil {
		return fmt.Errorf("invite user error: %w", err)
	}

	if storedInviteeID == uuid.Nil {
		return fmt.Errorf("invitee already in conversation or invite failed")
	}

	if err := s.notifications.InvalidateMembership(ctx, inviteeID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	if err := s.saveInvitedMessage(ctx, conversationID, inviteeID); err != nil {
		return fmt.Errorf("send invited message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	return nil
}

func (s *conversationService) Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error {
	isOwner, err := s.queries.IsMemberOwner(conversationID, kickerID)
	if err != nil {
		return fmt.Errorf("is owner error: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("user is not owner: %w", domain.ErrorUserNotOwner)
	}

	isMember, err := s.queries.IsMember(conversationID, targetID)
	if err != nil {
		return fmt.Errorf("is member error: %w", err)
	}
	if !isMember {
		return fmt.Errorf("target is not in conversation: %w", domain.ErrorUserNotInConversation)
	}

	rowsAffected, err := s.queries.KickParticipantAtomic(conversationID, targetID)
	if err != nil {
		return fmt.Errorf("kick participant error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("kick failed")
	}

	if err := s.saveLeftMessage(ctx, conversationID, targetID); err != nil {
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

func (s *conversationService) saveJoinedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeJoinedConversation, "")

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store joined message error: %w", err)
	}

	if !stored {
		return fmt.Errorf("system message validation failed")
	}

	messageDTO, err := s.queries.GetNotificationMessage(message.ID, userID)
	if err != nil {
		return fmt.Errorf("get message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) saveLeftMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeLeftConversation, "")

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store left message error: %w", err)
	}

	if !stored {
		return fmt.Errorf("system message validation failed")
	}

	messageDTO, err := s.queries.GetNotificationMessage(message.ID, userID)
	if err != nil {
		return fmt.Errorf("get message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) saveInvitedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeInvitedConversation, "")

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store invited message error: %w", err)
	}

	if !stored {
		return fmt.Errorf("system message validation failed")
	}

	messageDTO, err := s.queries.GetNotificationMessage(message.ID, userID)
	if err != nil {
		return fmt.Errorf("get message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}

func (s *conversationService) saveRenamedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, newName string) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeRenamedConversation, newName)

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store renamed message error: %w", err)
	}

	if !stored {
		return fmt.Errorf("system message validation failed")
	}

	messageDTO, err := s.queries.GetNotificationMessage(message.ID, userID)
	if err != nil {
		return fmt.Errorf("get message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}
