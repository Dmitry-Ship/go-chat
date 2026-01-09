package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"github.com/google/uuid"
)

type ParticipantCacheDecorator struct {
	cache CacheClient
	repo  domain.ParticipantRepository
}

func NewParticipantCacheDecorator(repo domain.ParticipantRepository, cache CacheClient) *ParticipantCacheDecorator {
	return &ParticipantCacheDecorator{
		cache: cache,
		repo:  repo,
	}
}

func (d *ParticipantCacheDecorator) GetByConversationIDAndUserID(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	key := ParticipantsKey(conversationID.String())

	data, err := d.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	if data != nil {
		var cachedData ParticipantsCache
		if err := json.Unmarshal(data, &cachedData); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		for _, p := range cachedData.Participants {
			if p.UserID == userID.String() {
				convID, _ := uuid.Parse(p.ConversationID)
				participantID, _ := uuid.Parse(p.ID)
				userIDParsed, _ := uuid.Parse(p.UserID)

				return &domain.Participant{
					ID:             participantID,
					ConversationID: convID,
					UserID:         userIDParsed,
				}, nil
			}
		}
	}

	participant, err := d.repo.GetByConversationIDAndUserID(ctx, conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("repo get by conversation id and user id error: %w", err)
	}

	return participant, nil
}

func (d *ParticipantCacheDecorator) GetIDsByConversationID(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	key := ParticipantsKey(conversationID.String())

	data, err := d.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	if data != nil {
		var cachedData ParticipantsCache
		if err := json.Unmarshal(data, &cachedData); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		ids := make([]uuid.UUID, len(cachedData.Participants))
		for i, p := range cachedData.Participants {
			ids[i], _ = uuid.Parse(p.UserID)
		}
		return ids, nil
	}

	ids, err := d.repo.GetIDsByConversationID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("repo get ids by conversation id error: %w", err)
	}

	participants := make([]*domain.Participant, len(ids))
	for i, id := range ids {
		participants[i] = &domain.Participant{
			UserID:         id,
			ConversationID: conversationID,
		}
	}

	data, err = SerializeParticipants(participants)
	if err != nil {
		return nil, fmt.Errorf("serialize participants error: %w", err)
	}

	if err := d.cache.Set(ctx, key, data, TTLParticipants); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return ids, nil
}

func (d *ParticipantCacheDecorator) Store(ctx context.Context, participant *domain.Participant) error {
	if err := d.repo.Store(ctx, participant); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateParticipantsCache(ctx, participant.ConversationID.String())
	d.invalidateUserConvListCache(ctx, participant.UserID.String())

	return nil
}

func (d *ParticipantCacheDecorator) Delete(ctx context.Context, participantID uuid.UUID) error {
	if err := d.repo.Delete(ctx, participantID); err != nil {
		return fmt.Errorf("repo delete error: %w", err)
	}

	d.invalidateUserConvListCache(ctx, participantID.String())

	return nil
}

func (d *ParticipantCacheDecorator) invalidateParticipantsCache(ctx context.Context, conversationID string) {
	_ = d.cache.Delete(ctx, ParticipantsKey(conversationID))
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

func (d *ParticipantCacheDecorator) invalidateUserConvListCache(ctx context.Context, userID string) {
	_ = d.cache.Delete(ctx, UserConvListKey(userID))
}
