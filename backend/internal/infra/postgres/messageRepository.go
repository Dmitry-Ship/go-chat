package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"
	"GitHub/go-chat/backend/internal/presentation"
	"GitHub/go-chat/backend/internal/readModel"

	"github.com/google/uuid"
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

func (r *messageRepository) StoreSystemMessage(ctx context.Context, message *domain.Message) (bool, error) {
	params := db.StoreSystemMessageParams{
		ID:             uuidToPgtype(message.ID),
		ConversationID: uuidToPgtype(message.ConversationID),
		UserID:         uuidToPgtype(message.UserID),
		Content:        message.Content.String(),
		Type:           int32(toMessageTypePersistence(message.Type)),
	}

	rowsAffected, err := r.queries.StoreSystemMessage(ctx, params)
	if err != nil {
		return false, fmt.Errorf("store system message error: %w", err)
	}

	return rowsAffected > 0, nil
}

func (r *messageRepository) StoreSystemMessageAndReturn(ctx context.Context, message *domain.Message, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	params := db.StoreSystemMessageParams{
		ID:             uuidToPgtype(message.ID),
		ConversationID: uuidToPgtype(message.ConversationID),
		UserID:         uuidToPgtype(message.UserID),
		Content:        message.Content.String(),
		Type:           int32(toMessageTypePersistence(message.Type)),
	}

	rowsAffected, err := r.queries.StoreSystemMessage(ctx, params)
	if err != nil {
		return readModel.MessageDTO{}, fmt.Errorf("store system message error: %w", err)
	}

	if rowsAffected == 0 {
		return readModel.MessageDTO{}, fmt.Errorf("system message validation failed")
	}

	dto, err := r.queries.GetNotificationMessageRaw(ctx, uuidToPgtype(message.ID))
	if err != nil {
		return readModel.MessageDTO{}, fmt.Errorf("get message error: %w", err)
	}

	formatter := presentation.NewMessageFormatter()
	rawMessage := readModel.RawMessageDTO{
		ID:             pgtypeToUUID(dto.ID),
		Type:           uint8(dto.Type),
		CreatedAt:      dto.CreatedAt.Time,
		ConversationID: pgtypeToUUID(dto.ConversationID),
		Content:        dto.Content,
		UserID:         pgtypeToUUID(dto.UserID),
		UserName:       dto.UserName.String,
		UserAvatar:     dto.UserAvatar.String,
	}

	return formatter.FormatMessageDTO(rawMessage, requestUserID), nil
}
