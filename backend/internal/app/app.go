package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ConversationService        services.ConversationService
	AuthService                services.AuthService
	NotificationClientRegister services.NotificationClientRegister
	MessagingService           services.MessagingService
}

type Queries struct {
	Users         readModel.UserQueryRepository
	Conversations readModel.ConversationQueryRepository
	Messages      readModel.MessageQueryRepository
}

func NewApp(db *gorm.DB, redisClient *redis.Client) *App {
	messagesRepository := postgres.NewMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	eventsPubSub := domain.NewPubsub()
	messagingService := services.NewMessagingService(messagesRepository, eventsPubSub)
	conversationService := services.NewConversationService(conversationsRepository, participantRepository, eventsPubSub)
	authService := services.NewAuthService(usersRepository)
	notificationService := services.NewNotificationService(redisClient, notificationTopicRepository)

	notificationEventHandlers := services.NewNotificationEventHandlers(eventsPubSub, notificationService, messagesRepository)
	messagesEventHandlers := services.NewMessagesEventHandlers(eventsPubSub, messagingService)

	go notificationEventHandlers.Run()
	go messagesEventHandlers.Run()
	go notificationService.Run()

	return &App{
		Commands: Commands{
			ConversationService:        conversationService,
			AuthService:                authService,
			NotificationClientRegister: notificationService,
			MessagingService:           messagingService,
		},
		Queries: Queries{
			Users:         usersRepository,
			Conversations: conversationsRepository,
			Messages:      messagesRepository,
		},
	}
}
