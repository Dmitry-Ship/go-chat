package app

import (
	"GitHub/go-chat/backend/pkg/postgres"
	"GitHub/go-chat/backend/pkg/readModel"
	"GitHub/go-chat/backend/pkg/services"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"gorm.io/gorm"
)

type App struct {
	Commands commands
	Queries  queries
}

type commands struct {
	ConversationService services.ConversationService
	AuthService         services.AuthService
}

type queries struct {
	UsersRepository        readModel.UserQueryRepository
	ConversationRepository readModel.ConversationQueryRepository
	MessageRepository      readModel.MessageQueryRepository
	ParticipantRepository  readModel.ParticipantQueryRepository
}

func NewApp(db *gorm.DB, hub ws.Hub) *App {
	messagesRepository := postgres.NewMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)

	conversationWSResolver := services.NewConversationWSResolver(participantRepository, messagesRepository, hub)

	return &App{
		Commands: commands{
			ConversationService: services.NewConversationService(conversationsRepository, participantRepository, messagesRepository, conversationWSResolver),
			AuthService:         services.NewAuthService(usersRepository, usersRepository),
		},
		Queries: queries{
			UsersRepository:        usersRepository,
			ConversationRepository: conversationsRepository,
			MessageRepository:      messagesRepository,
			ParticipantRepository:  participantRepository,
		},
	}
}
