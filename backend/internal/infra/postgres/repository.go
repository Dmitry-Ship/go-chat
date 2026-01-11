package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func newRepository(pool *pgxpool.Pool, queries *db.Queries) *repository {
	return &repository{
		pool:    pool,
		queries: queries,
	}
}

func (r *repository) withTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

func uuidToPgtype(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func pgtypeToUUID(u pgtype.UUID) uuid.UUID {
	return uuid.UUID(u.Bytes)
}

var conversationTypesMap = map[uint8]domain.ConversationType{
	0: domain.ConversationTypeGroup,
	1: domain.ConversationTypeDirect,
}

func toConversationTypePersistence(conversationType domain.ConversationType) uint8 {
	for k, v := range conversationTypesMap {
		if v == conversationType {
			return k
		}
	}
	return 0
}

var messageTypesMap = map[uint8]domain.MessageType{
	0: domain.MessageTypeUser,
	1: domain.MessageTypeSystem,
}

func toMessageTypePersistence(messageType domain.MessageType) uint8 {
	for k, v := range messageTypesMap {
		if v == messageType {
			return k
		}
	}
	return 0
}

func MessageTypePersistenceToDomain(persistenceType uint8) domain.MessageType {
	if messageType, ok := messageTypesMap[persistenceType]; ok {
		return messageType
	}
	return domain.MessageTypeUser
}
