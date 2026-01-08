package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	*repository
}

func NewMessageRepository(pool *pgxpool.Pool) *messageRepository {
	return &messageRepository{
		repository: newRepository(pool, db.New(pool)),
	}
}

func (r *messageRepository) Store(ctx context.Context, message *domain.Message) error {
	params := db.StoreMessageParams{
		ID:             uuidToPgtype(message.ID),
		ConversationID: uuidToPgtype(message.ConversationID),
		UserID:         uuidToPgtype(message.UserID),
		Content:        message.Content.String(),
		Type:           int32(toMessageTypePersistence(message.Type)),
	}

	if err := r.queries.StoreMessage(ctx, params); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}
