package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"

	"gorm.io/gorm"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ConversationService  services.ConversationService
	AuthService          services.AuthService
	NotificationsService services.NotificationsService
	MessagingService     services.MessagingService
}

type Queries struct {
	Users         readModel.UserQueryRepository
	Conversations readModel.ConversationQueryRepository
	Messages      readModel.MessageQueryRepository
}

func NewApp(db *gorm.DB, hub ws.Hub) *App {
	messagesRepository := postgres.NewMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	eventsPubSub := domain.NewPubsub()
	messagingService := services.NewMessagingService(messagesRepository, eventsPubSub)
	conversationService := services.NewConversationService(conversationsRepository, participantRepository, eventsPubSub)
	authService := services.NewAuthService(usersRepository)
	notificationsService := services.NewNotificationsService(notificationTopicRepository, hub)

	notificationEventHandlers := services.NewNotificationEventHandlers(eventsPubSub, notificationsService, messagesRepository)
	messagesEventHandlers := services.NewMessagesEventHandlers(eventsPubSub, messagingService)

	go notificationEventHandlers.Run()
	go messagesEventHandlers.Run()

	return &App{
		Commands: Commands{
			ConversationService:  conversationService,
			AuthService:          authService,
			NotificationsService: notificationsService,
			MessagingService:     messagingService,
		},
		Queries: Queries{
			Users:         usersRepository,
			Conversations: conversationsRepository,
			Messages:      messagesRepository,
		},
	}
}
