package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
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
	eventsPubSub := domain.NewPubsub()

	messagesRepository := postgres.NewMessageRepository(db, eventsPubSub)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db, eventsPubSub)
	participantRepository := postgres.NewParticipantRepository(db, eventsPubSub)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	messagingService := services.NewMessagingService(messagesRepository)
	conversationService := services.NewConversationService(conversationsRepository, participantRepository)
	authService := services.NewAuthService(usersRepository)
	notificationService := services.NewNotificationService(redisClient, notificationTopicRepository)

	notificationEventHandlers := domainEventsHandlers.NewNotificationEventHandlers(eventsPubSub, notificationService, messagesRepository)
	messagesEventHandlers := domainEventsHandlers.NewMessagesEventHandlers(eventsPubSub, messagingService)

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
