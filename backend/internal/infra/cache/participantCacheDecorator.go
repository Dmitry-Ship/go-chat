package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/repository"
	"github.com/google/uuid"
)

type ParticipantCacheDecorator struct {
	cache CacheClient
	repo  repository.ParticipantRepository
}

func NewParticipantCacheDecorator(repo repository.ParticipantRepository, cache CacheClient) *ParticipantCacheDecorator {
	return &ParticipantCacheDecorator{
		cache: cache,
		repo:  repo,
	}
}

func (d *ParticipantCacheDecorator) Store(ctx context.Context, participant *domain.Participant) error {
	if err := d.repo.Store(ctx, participant); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateParticipantsCache(ctx, participant.ConversationID.String())
	d.invalidateUserConvListCache(ctx, participant.UserID.String())

	return nil
}

func (d *ParticipantCacheDecorator) GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	key := UserConvListKey(userID.String())

	data, err := d.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	if data != nil {
		var cachedData []uuid.UUID
		if err := json.Unmarshal(data, &cachedData); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		return cachedData, nil
	}

	ids, err := d.repo.GetConversationIDsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo get conversation ids by user id error: %w", err)
	}

	data, err = json.Marshal(ids)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	if err := d.cache.Set(ctx, key, data, TTLUserConvList); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return ids, nil
}

func (d *ParticipantCacheDecorator) invalidateParticipantsCache(ctx context.Context, conversationID string) {
	_ = d.cache.Delete(ctx, ParticipantsKey(conversationID))
}

func (d *ParticipantCacheDecorator) invalidateUserConvListCache(ctx context.Context, userID string) {
	_ = d.cache.Delete(ctx, UserConvListKey(userID))
}
