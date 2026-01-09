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

type directConversationRepository struct {
	*repository
}

func NewDirectConversationRepository(pool *pgxpool.Pool) *directConversationRepository {
	return &directConversationRepository{
		repository: newRepository(pool, db.New(pool)),
	}
}

func (r *directConversationRepository) Store(ctx context.Context, conversation *domain.DirectConversation) error {
	return r.withTx(ctx, func(tx pgx.Tx) error {
		qtx := r.queries.WithTx(tx)

		conversationParams := db.StoreConversationParams{
			ID:       uuidToPgtype(conversation.ID),
			Type:     int32(toConversationTypePersistence(conversation.Type)),
			IsActive: conversation.IsActive,
		}

		if err := qtx.StoreConversation(ctx, conversationParams); err != nil {
			return fmt.Errorf("create conversation error: %w", err)
		}

		participantIDs := make([]pgtype.UUID, len(conversation.Participants))
		userIDs := make([]pgtype.UUID, len(conversation.Participants))
		for i, p := range conversation.Participants {
			participantIDs[i] = uuidToPgtype(p.ID)
			userIDs[i] = uuidToPgtype(p.UserID)
		}

		if err := qtx.StoreParticipantsBatch(ctx, db.StoreParticipantsBatchParams{
			Column1:        participantIDs,
			ConversationID: uuidToPgtype(conversation.ID),
			Column3:        userIDs,
		}); err != nil {
			return fmt.Errorf("create participants error: %w", err)
		}

		return nil
	})
}

func (r *directConversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DirectConversation, error) {
	conv, err := r.queries.GetDirectConversationWithParticipants(ctx, uuidToPgtype(id))
	if err != nil {
		return nil, fmt.Errorf("get direct conversation error: %w", err)
	}

	// Parse the ARRAY_AGG result - PostgreSQL returns this as a slice of [16]byte
	participantUserIDs, ok := conv.ParticipantUserIds.([]interface{})
	if !ok || len(participantUserIDs) == 0 {
		return nil, fmt.Errorf("direct conversation not found")
	}

	participantsDomain := make([]domain.Participant, len(participantUserIDs))
	for i, userIDRaw := range participantUserIDs {
		// Convert from [16]byte to uuid.UUID
		userIDBytes, ok := userIDRaw.([16]byte)
		if !ok {
			return nil, fmt.Errorf("invalid user id format")
		}
		userID := uuid.UUID(userIDBytes)
		participantsDomain[i] = domain.Participant{
			UserID:         userID,
			ConversationID: pgtypeToUUID(conv.ID),
			IsActive:       true,
		}
	}

	return &domain.DirectConversation{
		Participants: participantsDomain,
		Conversation: domain.Conversation{
			ID:       pgtypeToUUID(conv.ID),
			Type:     conversationTypesMap[uint8(conv.Type)],
			IsActive: conv.IsActive,
		},
	}, nil
}

func (r *directConversationRepository) GetID(ctx context.Context, firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	conv, err := r.queries.GetDirectConversationBetweenUsers(ctx, db.GetDirectConversationBetweenUsersParams{
		UserID:   uuidToPgtype(firstUserID),
		UserID_2: uuidToPgtype(secondUserID),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("get conversation id error: %w", err)
	}

	return pgtypeToUUID(conv.ID), nil
}
