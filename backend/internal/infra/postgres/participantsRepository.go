package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type participantRepository struct {
	*repository
}

func NewParticipantRepository(pool *pgxpool.Pool) *participantRepository {
	return &participantRepository{
		repository: newRepository(pool, db.New(pool)),
	}
}

func (r *participantRepository) Store(ctx context.Context, participant *domain.Participant) error {
	params := db.StoreParticipantParams{
		ID:             uuidToPgtype(participant.ID),
		ConversationID: uuidToPgtype(participant.ConversationID),
		UserID:         uuidToPgtype(participant.UserID),
	}

	if err := r.queries.StoreParticipant(ctx, params); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	return nil
}

func (r *participantRepository) GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	conversationIDs, err := r.queries.GetConversationIDsByUserID(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("get user conversations error: %w", err)
	}

	ids := make([]uuid.UUID, len(conversationIDs))
	for i, id := range conversationIDs {
		ids[i] = pgtypeToUUID(id)
	}

	return ids, nil
}
