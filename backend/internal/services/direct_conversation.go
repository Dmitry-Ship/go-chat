package services

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type DirectConversationService interface {
	StartDirectConversation(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error)
}

type directConversationService struct {
	directConversations domain.DirectConversationRepository
	notifications       NotificationService
	cache               CacheService
}

func NewDirectConversationService(
	directConversations domain.DirectConversationRepository,
	notifications NotificationService,
	cache CacheService,
) DirectConversationService {
	return &directConversationService{
		directConversations: directConversations,
		notifications:       notifications,
		cache:               cache,
	}
}

func (s *directConversationService) StartDirectConversation(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error) {
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
