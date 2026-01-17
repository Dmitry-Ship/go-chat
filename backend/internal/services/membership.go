package services

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/repository"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type MembershipService interface {
	Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
	Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error
}

type membershipService struct {
	participants  repository.ParticipantRepository
	queries       readModel.QueriesRepository
	messages      MessageService
	notifications NotificationService
	cache         CacheService
}

func NewMembershipService(
	participants repository.ParticipantRepository,
	queries readModel.QueriesRepository,
	messages MessageService,
	notifications NotificationService,
	cache CacheService,
) MembershipService {
	return &membershipService{
		participants:  participants,
		queries:       queries,
		messages:      messages,
		notifications: notifications,
		cache:         cache,
	}
}

func (s *membershipService) Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	participant := domain.NewParticipant(uuid.New(), conversationID, userID)

	if err := s.participants.Store(ctx, participant); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	joinedMessage, err := domain.NewMessage(conversationID, userID, domain.MessageTypeSystem, "joined the conversation")
	if err != nil {
		return fmt.Errorf("create joined message error: %w", err)
	}

	if _, err := s.messages.Send(ctx, joinedMessage); err != nil {
		return fmt.Errorf("store joined message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, userID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *membershipService) Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
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

	leftMessage, err := domain.NewMessage(conversationID, userID, domain.MessageTypeSystem, "left the conversation")
	if err != nil {
		return fmt.Errorf("create left message error: %w", err)
	}

	if _, err := s.messages.Send(ctx, leftMessage); err != nil {
		return fmt.Errorf("store left message error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, userID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	return nil
}

func (s *membershipService) Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error {
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

	invitedMessage, err := domain.NewMessage(conversationID, inviteeID, domain.MessageTypeSystem, "was invited to the conversation")
	if err != nil {
		return fmt.Errorf("create invited message error: %w", err)
	}

	if _, err := s.messages.Send(ctx, invitedMessage); err != nil {
		return fmt.Errorf("store invited message error: %w", err)
	}

	if err := s.notifications.InvalidateMembership(ctx, inviteeID); err != nil {
		return fmt.Errorf("invalidate membership error: %w", err)
	}

	if err := s.cache.InvalidateParticipants(ctx, conversationID); err != nil {
		return fmt.Errorf("invalidate cache error: %w", err)
	}

	return nil
}

func (s *membershipService) Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error {
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

	rowsAffected, err := s.queries.LeaveConversationAtomic(conversationID, targetID)
	if err != nil {
		return fmt.Errorf("kick participant error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("kick failed")
	}

	kickedMessage, err := domain.NewMessage(conversationID, targetID, domain.MessageTypeSystem, "was kicked from the conversation")
	if err != nil {
		return fmt.Errorf("create kicked message error: %w", err)
	}

	if _, err := s.messages.Send(ctx, kickedMessage); err != nil {
		return fmt.Errorf("store left message error: %w", err)
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
