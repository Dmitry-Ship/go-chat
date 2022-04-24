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
	ConversationService  services.ConversationService
	AuthService          services.AuthService
	NotificationsService services.NotificationsService
	MessagingService     services.MessagingService
}

type queries struct {
	UsersRepository        readModel.UserQueryRepository
	ConversationRepository readModel.ConversationQueryRepository
	MessageRepository      readModel.MessageQueryRepository
}

func NewApp(db *gorm.DB, hub ws.Hub) *App {
	messagesRepository := postgres.NewMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	notificationsService := services.NewNotificationsService(messagesRepository, notificationTopicRepository, hub)
	messagingService := services.NewMessagingService(messagesRepository, notificationsService)
	conversationService := services.NewConversationService(conversationsRepository, participantRepository, messagingService, notificationsService)
	authService := services.NewAuthService(usersRepository)

	return &App{
		Commands: commands{
			ConversationService:  conversationService,
			AuthService:          authService,
			NotificationsService: notificationsService,
			MessagingService:     messagingService,
		},
		Queries: queries{
			UsersRepository:        usersRepository,
			ConversationRepository: conversationsRepository,
			MessageRepository:      messagesRepository,
		},
	}
}
