package app

import (
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	"GitHub/go-chat/backend/internal/services"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"

	"gorm.io/gorm"
)

type Commands struct {
	ConversationService services.ConversationService
	AuthService         services.AuthService
	ClientRegister      services.ClientRegister
}

func NewCommands(ctx context.Context, eventPublisher infra.EventPublisher, db *gorm.DB, activeClients ws.ActiveClients) *Commands {
	messagesRepository := postgres.NewMessageRepository(db, eventPublisher)
	usersRepository := postgres.NewUserRepository(db)
	publicConversationsRepository := postgres.NewPublicConversationRepository(db, eventPublisher)
	privateConversationsRepository := postgres.NewPrivateConversationRepository(db, eventPublisher)
	participantRepository := postgres.NewParticipantRepository(db, eventPublisher)
	jwtTokens := services.NewJWTokens()

	return &Commands{
		ConversationService: services.NewConversationService(publicConversationsRepository, privateConversationsRepository, participantRepository, messagesRepository),
		AuthService:         services.NewAuthService(usersRepository, jwtTokens),
		ClientRegister:      services.NewClientRegister(activeClients),
	}
}
