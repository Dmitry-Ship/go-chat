package repository

import (
	"context"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"

	"github.com/google/uuid"
)

type UserRepository interface {
	Store(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

type MessageRepository interface {
	Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error)
}

type ParticipantRepository interface {
	Store(ctx context.Context, participant *domain.Participant) error
	Delete(ctx context.Context, participantID uuid.UUID) error
	GetByConversationIDAndUserID(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*domain.Participant, error)
	GetIDsByConversationID(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
	GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type DirectConversationRepository interface {
	Store(ctx context.Context, conversation *domain.DirectConversation) error
	GetID(ctx context.Context, firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.DirectConversation, error)
}

type GroupConversationRepository interface {
	Store(ctx context.Context, conversation *domain.GroupConversation) error
	Update(ctx context.Context, conversation *domain.GroupConversation) error
	Rename(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GroupConversation, error)
}
