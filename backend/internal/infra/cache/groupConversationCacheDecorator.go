package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"github.com/google/uuid"
)

type GroupConversationCacheDecorator struct {
	cache CacheClient
	repo  domain.GroupConversationRepository
}

func NewGroupConversationCacheDecorator(repo domain.GroupConversationRepository, cache CacheClient) *GroupConversationCacheDecorator {
	return &GroupConversationCacheDecorator{
		cache: cache,
		repo:  repo,
	}
}

func (d *GroupConversationCacheDecorator) GetByID(ctx context.Context, id uuid.UUID) (*domain.GroupConversation, error) {
	key := ConversationKey(id.String())

	data, err := d.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	if data != nil {
		var conv domain.GroupConversation
		if err := json.Unmarshal(data, &conv); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}
		return &conv, nil
	}

	conv, err := d.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo get by id error: %w", err)
	}

	data, err = json.Marshal(conv)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	if err := d.cache.Set(ctx, key, data, TTLConversation); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return conv, nil
}

func (d *GroupConversationCacheDecorator) Store(ctx context.Context, conversation *domain.GroupConversation) error {
	if err := d.repo.Store(ctx, conversation); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateConversationCache(ctx, conversation.ID.String())

	return nil
}

func (d *GroupConversationCacheDecorator) Update(ctx context.Context, conversation *domain.GroupConversation) error {
	if err := d.repo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("repo update error: %w", err)
	}

	d.invalidateConversationCache(ctx, conversation.ID.String())

	return nil
}

func (d *GroupConversationCacheDecorator) invalidateConversationCache(ctx context.Context, conversationID string) {
	_ = d.cache.Delete(ctx, ConversationKey(conversationID))
	_ = d.cache.Delete(ctx, ConvMetaKey(conversationID))
	_ = d.cache.DeletePattern(ctx, ParticipantsKey(conversationID))
}
