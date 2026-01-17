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

type GroupConversationService interface {
	CreateGroupConversation(ctx context.Context, conversationID uuid.UUID, name string, userID uuid.UUID) error
	DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Rename(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, name string) error
}

type groupConversationService struct {
	groupConversations repository.GroupConversationRepository
	queries            readModel.QueriesRepository
	messages           MessageService
	notifications      NotificationService
	cache              CacheService
}

func NewGroupConversationService(
	groupConversations repository.GroupConversationRepository,
	queries readModel.QueriesRepository,
	messages MessageService,
	notifications NotificationService,
	cache CacheService,
) GroupConversationService {
	return &groupConversationService{
		groupConversations: groupConversations,
		queries:            queries,
		messages:           messages,
		notifications:      notifications,
		cache:              cache,
	}
}

func (s *groupConversationService) CreateGroupConversation(ctx context.Context, conversationID uuid.UUID, name string, userID uuid.UUID) error {
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

func (s *groupConversationService) DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
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

func (s *groupConversationService) Rename(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, name string) error {
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

	renameMessage, err := domain.NewMessage(conversationID, userID, domain.MessageTypeSystem, fmt.Sprintf("renamed the conversation to %s", name))
	if err != nil {
		return fmt.Errorf("create rename message error: %w", err)
	}

	if _, err := s.messages.Send(ctx, renameMessage); err != nil {
		return fmt.Errorf("store renamed message error: %w", err)
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
