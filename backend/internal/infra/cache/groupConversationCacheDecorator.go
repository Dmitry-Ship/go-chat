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

func (d *GroupConversationCacheDecorator) GetByID(id uuid.UUID) (*domain.GroupConversation, error) {
	key := ConversationKey(id.String())

	data, err := d.cache.Get(context.Background(), key)
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

	conv, err := d.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("repo get by id error: %w", err)
	}

	data, err = json.Marshal(conv)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	if err := d.cache.Set(context.Background(), key, data, TTLConversation); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return conv, nil
}

func (d *GroupConversationCacheDecorator) Store(conversation *domain.GroupConversation) error {
	if err := d.repo.Store(conversation); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateConversationCache(conversation.ID.String())

	return nil
}

func (d *GroupConversationCacheDecorator) Update(conversation *domain.GroupConversation) error {
	if err := d.repo.Update(conversation); err != nil {
		return fmt.Errorf("repo update error: %w", err)
	}

	d.invalidateConversationCache(conversation.ID.String())

	return nil
}

func (d *GroupConversationCacheDecorator) invalidateConversationCache(conversationID string) {
	_ = d.cache.Delete(context.Background(), ConversationKey(conversationID))
	_ = d.cache.Delete(context.Background(), ConvMetaKey(conversationID))
	_ = d.cache.DeletePattern(context.Background(), ParticipantsKey(conversationID))
}
