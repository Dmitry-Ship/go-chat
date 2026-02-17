package cache

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/repository"
	"github.com/google/uuid"
)

type GroupConversationCacheDecorator struct {
	cache CacheClient
	repo  repository.GroupConversationRepository
}

func NewGroupConversationCacheDecorator(repo repository.GroupConversationRepository, cache CacheClient) *GroupConversationCacheDecorator {
	return &GroupConversationCacheDecorator{
		cache: cache,
		repo:  repo,
	}
}

func (d *GroupConversationCacheDecorator) Store(ctx context.Context, conversation *domain.GroupConversation) error {
	if err := d.repo.Store(ctx, conversation); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateConversationCache(ctx, conversation.ID.String())

	return nil
}

func (d *GroupConversationCacheDecorator) Rename(ctx context.Context, id uuid.UUID, name string) error {
	if err := d.repo.Rename(ctx, id, name); err != nil {
		return fmt.Errorf("repo rename error: %w", err)
	}

	d.invalidateConversationCache(ctx, id.String())

	return nil
}

func (d *GroupConversationCacheDecorator) Delete(ctx context.Context, id uuid.UUID) error {
	if err := d.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo delete error: %w", err)
	}

	d.invalidateConversationCache(ctx, id.String())

	return nil
}

func (d *GroupConversationCacheDecorator) invalidateConversationCache(ctx context.Context, conversationID string) {
	_ = d.cache.Delete(ctx, ConversationKey(conversationID))
	_ = d.cache.Delete(ctx, ConvMetaKey(conversationID))
	_ = d.cache.DeletePattern(ctx, ParticipantsKey(conversationID))
}
