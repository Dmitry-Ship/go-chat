package postgres

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"
	"GitHub/go-chat/backend/internal/presentation"
	"GitHub/go-chat/backend/internal/readModel"

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

func (r *messageRepository) Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error) {
	params := db.StoreMessageAndReturnParams{
		ID:             uuidToPgtype(message.ID),
		ConversationID: uuidToPgtype(message.ConversationID),
		UserID:         uuidToPgtype(message.UserID),
		Content:        message.Content.String(),
		Type:           int32(toMessageTypePersistence(message.Type)),
	}

	msg, err := r.queries.StoreMessageAndReturn(ctx, params)
	if err != nil {
		return readModel.MessageDTO{}, fmt.Errorf("store message error: %w", err)
	}

	formatter := presentation.NewMessageFormatter()
	rawMessage := readModel.RawMessageDTO{
		ID:             pgtypeToUUID(msg.ID),
		Type:           uint8(msg.Type),
		CreatedAt:      msg.CreatedAt.Time,
		ConversationID: pgtypeToUUID(msg.ConversationID),
		Content:        msg.FormattedText,
		UserID:         pgtypeToUUID(msg.UserID),
		UserName:       msg.UserName,
		UserAvatar:     msg.UserAvatar.String,
	}

	return formatter.FormatMessageDTO(rawMessage), nil
}
