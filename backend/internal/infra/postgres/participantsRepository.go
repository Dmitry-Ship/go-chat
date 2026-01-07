package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type participantRepository struct {
	*repository
}

func NewParticipantRepository(pool *pgxpool.Pool, eventPublisher *infra.EventBus) *participantRepository {
	return &participantRepository{
		repository: newRepository(pool, db.New(pool), eventPublisher),
	}
}

func (r *participantRepository) Store(participant *domain.Participant) error {
	ctx := context.Background()
	params := db.StoreParticipantParams{
		ID:             uuidToPgtype(participant.ID),
		ConversationID: uuidToPgtype(participant.ConversationID),
		UserID:         uuidToPgtype(participant.UserID),
		IsActive:       participant.IsActive,
	}

	if err := r.queries.StoreParticipant(ctx, params); err != nil {
		return fmt.Errorf("store participant error: %w", err)
	}

	r.dispatchEvents(participant)
	return nil
}

func (r *participantRepository) Update(participant *domain.Participant) error {
	ctx := context.Background()
	params := db.UpdateParticipantParams{
		ID:       uuidToPgtype(participant.ID),
		IsActive: participant.IsActive,
	}

	if err := r.queries.UpdateParticipant(ctx, params); err != nil {
		return fmt.Errorf("update participant error: %w", err)
	}

	r.dispatchEvents(participant)
	return nil
}

func (r *participantRepository) GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	ctx := context.Background()
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

func (r *participantRepository) GetIDsByConversationID(conversationID uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
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

func (r *participantRepository) GetConversationIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	participants, err := r.queries.GetParticipantByID(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("get user conversations error: %w", err)
	}

	return []uuid.UUID{pgtypeToUUID(participants.ConversationID)}, nil
}
