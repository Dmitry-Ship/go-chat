package services

import (
	"context"

	"GitHub/go-chat/backend/internal/infra/cache"

	"github.com/google/uuid"
)

type CacheService interface {
	InvalidateConversation(ctx context.Context, conversationID uuid.UUID) error
	InvalidateParticipants(ctx context.Context, conversationID uuid.UUID) error
	InvalidateUserConversations(ctx context.Context, userID uuid.UUID) error
}

type cacheService struct {
	cache cache.CacheClient
}

func NewCacheService(cacheClient cache.CacheClient) CacheService {
	return &cacheService{
		cache: cacheClient,
	}
}

func (s *cacheService) InvalidateConversation(ctx context.Context, conversationID uuid.UUID) error {
	if err := s.cache.Delete(ctx, cache.ConversationKey(conversationID.String())); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, cache.ConvMetaKey(conversationID.String())); err != nil {
		return err
	}
	if err := s.cache.DeletePattern(ctx, cache.ParticipantsKey(conversationID.String())); err != nil {
		return err
	}

	return nil
}

func (s *cacheService) InvalidateParticipants(ctx context.Context, conversationID uuid.UUID) error {
	if err := s.cache.Delete(ctx, cache.ParticipantsKey(conversationID.String())); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, cache.ConvMetaKey(conversationID.String())); err != nil {
		return err
	}

	return nil
}

func (s *cacheService) InvalidateUserConversations(ctx context.Context, userID uuid.UUID) error {
	if err := s.cache.Delete(ctx, cache.UserConvListKey(userID.String())); err != nil {
		return err
	}

	return nil
}
