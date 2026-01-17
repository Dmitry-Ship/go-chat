package cache

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/repository"
	"github.com/google/uuid"
)

type UserCacheDecorator struct {
	cache CacheClient
	repo  repository.UserRepository
}

func NewUserCacheDecorator(repo repository.UserRepository, cache CacheClient) *UserCacheDecorator {
	return &UserCacheDecorator{
		cache: cache,
		repo:  repo,
	}
}

func (d *UserCacheDecorator) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	key := UserKey(id.String())

	data, err := d.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	if data != nil {
		cachedUser, err := DeserializeUser(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize user error: %w", err)
		}

		userID, err := uuid.Parse(cachedUser.ID)
		if err != nil {
			return nil, fmt.Errorf("parse user id error: %w", err)
		}

		return &domain.User{
			ID:           userID,
			Name:         cachedUser.Name,
			PasswordHash: "cached_dummy_password_hash",
			Avatar:       cachedUser.Avatar,
		}, nil
	}

	user, err := d.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo get by id error: %w", err)
	}

	data, err = SerializeUser(user)
	if err != nil {
		return nil, fmt.Errorf("serialize user error: %w", err)
	}

	if err := d.cache.Set(ctx, key, data, TTLUser); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return user, nil
}

func (d *UserCacheDecorator) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	key := UsernameKey(username)

	data, err := d.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	if data != nil {
		cachedUser, err := DeserializeUser(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize user error: %w", err)
		}

		userID, err := uuid.Parse(cachedUser.ID)
		if err != nil {
			return nil, fmt.Errorf("parse user id error: %w", err)
		}

		return &domain.User{
			ID:           userID,
			Name:         cachedUser.Name,
			PasswordHash: "cached_dummy_password_hash",
			Avatar:       cachedUser.Avatar,
		}, nil
	}

	user, err := d.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("repo find by username error: %w", err)
	}

	data, err = SerializeUser(user)
	if err != nil {
		return nil, fmt.Errorf("serialize user error: %w", err)
	}

	if err := d.cache.Set(ctx, key, data, TTLUser); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return user, nil
}

func (d *UserCacheDecorator) Store(ctx context.Context, user *domain.User) error {
	if err := d.repo.Store(ctx, user); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateUserCache(ctx, user)

	return nil
}

func (d *UserCacheDecorator) Update(ctx context.Context, user *domain.User) error {
	if err := d.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("repo update error: %w", err)
	}

	d.invalidateUserCache(ctx, user)

	return nil
}

func (d *UserCacheDecorator) invalidateUserCache(ctx context.Context, user *domain.User) {
	_ = d.cache.Delete(ctx, UserKey(user.ID.String()))
	_ = d.cache.Delete(ctx, UsernameKey(user.Name))
}
