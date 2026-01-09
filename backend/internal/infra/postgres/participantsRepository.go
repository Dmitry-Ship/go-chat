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
		IsActive:       participant.IsActive,
	}

	if err := r.queries.StoreParticipant(ctx, params); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	return nil
}

func (r *participantRepository) Update(ctx context.Context, participant *domain.Participant) error {
	params := db.UpdateParticipantParams{
		ID:       uuidToPgtype(participant.ID),
		IsActive: participant.IsActive,
	}

	if err := r.queries.UpdateParticipant(ctx, params); err != nil {
		return fmt.Errorf("update participant error: %w", err)
	}

	return nil
}

func (r *participantRepository) GetByConversationIDAndUserID(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	params := db.FindParticipantByConversationAndUserParams{
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(userID),
	}

	participant, err := r.queries.FindParticipantByConversationAndUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get participant error: %w", err)
	}

	return &domain.Participant{
		ID:             pgtypeToUUID(participant.ID),
		ConversationID: pgtypeToUUID(participant.ConversationID),
		UserID:         pgtypeToUUID(participant.UserID),
		IsActive:       participant.IsActive,
	}, nil
}

func (r *participantRepository) GetIDsByConversationID(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	participants, err := r.queries.GetParticipantsIDsByConversationID(ctx, uuidToPgtype(conversationID))

	if err != nil {
		return nil, fmt.Errorf("get participants error: %w", err)
	}

	ids := make([]uuid.UUID, len(participants))
	for i, p := range participants {
		ids[i] = pgtypeToUUID(p)
	}

	return ids, nil
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
