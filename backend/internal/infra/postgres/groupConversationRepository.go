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

type groupConversationRepository struct {
	*repository
}

func NewGroupConversationRepository(pool *pgxpool.Pool) *groupConversationRepository {
	return &groupConversationRepository{
		repository: newRepository(pool, db.New(pool)),
	}
}

func (r *groupConversationRepository) Store(ctx context.Context, conversation *domain.GroupConversation) error {
	return r.withTx(ctx, func(tx pgx.Tx) error {
		qtx := r.queries.WithTx(tx)

		conversationParams := db.StoreConversationParams{
			ID:   uuidToPgtype(conversation.Conversation.ID),
			Type: int32(toConversationTypePersistence(conversation.Type)),
		}

		if err := qtx.StoreConversation(ctx, conversationParams); err != nil {
			return fmt.Errorf("create conversation error: %w", err)
		}

		groupConversationParams := db.StoreGroupConversationParams{
			ID:             uuidToPgtype(conversation.ID),
			Name:           conversation.Name,
			Avatar:         pgtype.Text{String: conversation.Avatar, Valid: conversation.Avatar != ""},
			ConversationID: uuidToPgtype(conversation.Conversation.ID),
			OwnerID:        uuidToPgtype(conversation.Owner.UserID),
		}

		if err := qtx.StoreGroupConversation(ctx, groupConversationParams); err != nil {
			return fmt.Errorf("create group conversation error: %w", err)
		}

		participantParams := db.StoreParticipantParams{
			ID:             uuidToPgtype(conversation.Owner.ID),
			ConversationID: uuidToPgtype(conversation.Owner.ConversationID),
			UserID:         uuidToPgtype(conversation.Owner.UserID),
		}

		if err := qtx.StoreParticipant(ctx, participantParams); err != nil {
			return fmt.Errorf("create participant error: %w", err)
		}

		return nil
	})
}

func (r *groupConversationRepository) Rename(ctx context.Context, id uuid.UUID, name string) error {
	params := db.RenameGroupConversationParams{
		ConversationID: uuidToPgtype(id),
		Name:           name,
	}

	if err := r.queries.RenameGroupConversation(ctx, params); err != nil {
		return fmt.Errorf("rename group conversation error: %w", err)
	}

	return nil
}

func (r *groupConversationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteConversation(ctx, uuidToPgtype(id)); err != nil {
		return fmt.Errorf("delete conversation error: %w", err)
	}

	return nil
}
