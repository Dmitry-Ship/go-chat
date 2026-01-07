package cache

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"github.com/google/uuid"
)

type UserCacheDecorator struct {
	cache CacheClient
	repo  domain.UserRepository
}

func NewUserCacheDecorator(repo domain.UserRepository, cache CacheClient) *UserCacheDecorator {
	return &UserCacheDecorator{
		cache: cache,
		repo:  repo,
	}
}

func (d *UserCacheDecorator) GetByID(id uuid.UUID) (*domain.User, error) {
	key := UserKey(id.String())

	data, err := d.cache.Get(context.Background(), key)
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

	user, err := d.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("repo get by id error: %w", err)
	}

	data, err = SerializeUser(user)
	if err != nil {
		return nil, fmt.Errorf("serialize user error: %w", err)
	}

	if err := d.cache.Set(context.Background(), key, data, TTLUser); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return user, nil
}

func (d *UserCacheDecorator) FindByUsername(username string) (*domain.User, error) {
	key := UsernameKey(username)

	data, err := d.cache.Get(context.Background(), key)
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

	user, err := d.repo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("repo find by username error: %w", err)
	}

	data, err = SerializeUser(user)
	if err != nil {
		return nil, fmt.Errorf("serialize user error: %w", err)
	}

	if err := d.cache.Set(context.Background(), key, data, TTLUser); err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return user, nil
}

func (d *UserCacheDecorator) Store(user *domain.User) error {
	if err := d.repo.Store(user); err != nil {
		return fmt.Errorf("repo store error: %w", err)
	}

	d.invalidateUserCache(user)

	return nil
}

func (d *UserCacheDecorator) Update(user *domain.User) error {
	if err := d.repo.Update(user); err != nil {
		return fmt.Errorf("repo update error: %w", err)
	}

	d.invalidateUserCache(user)

	return nil
}

func (d *UserCacheDecorator) invalidateUserCache(user *domain.User) {
	_ = d.cache.Delete(context.Background(), UserKey(user.ID.String()))
	_ = d.cache.Delete(context.Background(), UsernameKey(user.Name))
}
