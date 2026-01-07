package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	*repository
}

func NewMessageRepository(pool *pgxpool.Pool, eventPublisher *infra.EventBus) *messageRepository {
	return &messageRepository{
		repository: newRepository(pool, db.New(pool), eventPublisher),
	}
}

func (r *messageRepository) Store(message *domain.Message) error {
	ctx := context.Background()
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

	r.dispatchEvents(message)
	return nil
}
