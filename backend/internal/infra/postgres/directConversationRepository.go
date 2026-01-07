package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type directConversationRepository struct {
	*repository
}

func NewDirectConversationRepository(pool *pgxpool.Pool, eventPublisher *infra.EventBus) *directConversationRepository {
	return &directConversationRepository{
		repository: newRepository(pool, db.New(pool), eventPublisher),
	}
}

func (r *directConversationRepository) Store(conversation *domain.DirectConversation) error {
	ctx := context.Background()
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

		for _, participant := range conversation.Participants {
			participantParams := db.StoreParticipantParams{
				ID:             uuidToPgtype(participant.ID),
				ConversationID: uuidToPgtype(participant.ConversationID),
				UserID:         uuidToPgtype(participant.UserID),
				IsActive:       participant.IsActive,
			}

			if err := qtx.StoreParticipant(ctx, participantParams); err != nil {
				return fmt.Errorf("create participant error: %w", err)
			}
		}

		return nil
	})
}

func (r *directConversationRepository) GetByID(id uuid.UUID) (*domain.DirectConversation, error) {
	ctx := context.Background()
	participants, err := r.queries.GetParticipantsIDsByConversationID(ctx, uuidToPgtype(id))
	if err != nil {
		return nil, fmt.Errorf("get direct conversation error: %w", err)
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("direct conversation not found")
	}

	conv, err := r.queries.GetConversationByID(ctx, uuidToPgtype(id))
	if err != nil {
		return nil, fmt.Errorf("get conversation error: %w", err)
	}

	participantsDomain := make([]domain.Participant, len(participants))
	for i, userID := range participants {
		participantsDomain[i] = domain.Participant{
			UserID:         pgtypeToUUID(userID),
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

func (r *directConversationRepository) GetID(firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	ctx := context.Background()
	conv, err := r.queries.GetDirectConversationBetweenUsers(ctx, db.GetDirectConversationBetweenUsersParams{
		UserID:   uuidToPgtype(firstUserID),
		UserID_2: uuidToPgtype(secondUserID),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("get conversation id error: %w", err)
	}

	return pgtypeToUUID(conv.ID), nil
}
