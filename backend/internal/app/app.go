package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres"
	"GitHub/go-chat/backend/internal/services"
	"context"

	"gorm.io/gorm"
)

type Commands struct {
	ConversationService services.ConversationService
	AuthService         services.AuthService
	MessagingService    services.MessagingService
}

func NewCommands(ctx context.Context, eventsPubSub domain.EventPublisher, db *gorm.DB) *Commands {
	messagesRepository := postgres.NewMessageRepository(db, eventsPubSub)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db, eventsPubSub)
	participantRepository := postgres.NewParticipantRepository(db, eventsPubSub)

	return &Commands{
		ConversationService: services.NewConversationService(conversationsRepository, participantRepository),
		AuthService:         services.NewAuthService(usersRepository),
		MessagingService:    services.NewMessagingService(messagesRepository),
	}
}
