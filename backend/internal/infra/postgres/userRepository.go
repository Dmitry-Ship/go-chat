package postgres

import (
	"context"
	"fmt"
	"time"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	*repository
}

func NewUserRepository(pool *pgxpool.Pool, eventPublisher *infra.EventBus) *userRepository {
	return &userRepository{
		repository: newRepository(pool, db.New(pool), eventPublisher),
	}
}

func (r *userRepository) Store(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	params := db.StoreUserParams{
		ID:           uuidToPgtype(user.ID),
		Avatar:       pgtype.Text{String: user.Avatar, Valid: user.Avatar != ""},
		Name:         user.Name,
		Password:     user.PasswordHash,
		RefreshToken: pgtype.Text{String: user.RefreshToken, Valid: user.RefreshToken != ""},
	}

	if err := r.queries.StoreUser(ctx, params); err != nil {
		return fmt.Errorf("store user error: %w", err)
	}

	r.dispatchEvents(user)
	return nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	params := db.UpdateUserParams{
		ID:           uuidToPgtype(user.ID),
		Avatar:       pgtype.Text{String: user.Avatar, Valid: user.Avatar != ""},
		Name:         user.Name,
		Password:     user.PasswordHash,
		RefreshToken: pgtype.Text{String: user.RefreshToken, Valid: user.RefreshToken != ""},
	}

	if err := r.queries.UpdateUser(ctx, params); err != nil {
		return fmt.Errorf("update user error: %w", err)
	}

	r.dispatchEvents(user)
	return nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := r.queries.GetUserByID(ctx, uuidToPgtype(id))
	if err != nil {
		return nil, fmt.Errorf("get user by id error: %w", err)
	}

	return &domain.User{
		ID:           pgtypeToUUID(user.ID),
		Avatar:       user.Avatar.String,
		Name:         user.Name,
		PasswordHash: user.Password,
		RefreshToken: user.RefreshToken.String,
	}, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := r.queries.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("find user by username error: %w", err)
	}

	return &domain.User{
		ID:           pgtypeToUUID(user.ID),
		Avatar:       user.Avatar.String,
		Name:         user.Name,
		PasswordHash: user.Password,
		RefreshToken: user.RefreshToken.String,
	}, nil
}

func (r *userRepository) UpdateRefreshToken(id uuid.UUID, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	params := db.UpdateUserRefreshTokenParams{
		ID:           uuidToPgtype(id),
		RefreshToken: pgtype.Text{String: refreshToken, Valid: refreshToken != ""},
	}

	if err := r.queries.UpdateUserRefreshToken(ctx, params); err != nil {
		return fmt.Errorf("update refresh token error: %w", err)
	}

	return nil
}
