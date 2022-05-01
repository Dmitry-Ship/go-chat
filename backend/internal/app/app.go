package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ConversationService services.ConversationService
	AuthService         services.AuthService
	NotificationService services.NotificationService
	MessagingService    services.MessagingService
}

type Queries struct {
	Users         readModel.UserQueryRepository
	Conversations readModel.ConversationQueryRepository
	Messages      readModel.MessageQueryRepository
}

func NewApp(ctx context.Context, eventsPubSub domain.EventPublisher, db *gorm.DB, redisClient *redis.Client) *App {
	messagesRepository := postgres.NewMessageRepository(db, eventsPubSub)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db, eventsPubSub)
	participantRepository := postgres.NewParticipantRepository(db, eventsPubSub)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	messagingService := services.NewMessagingService(messagesRepository)
	conversationService := services.NewConversationService(conversationsRepository, participantRepository)
	authService := services.NewAuthService(usersRepository)
	notificationService := services.NewNotificationService(ctx, redisClient, notificationTopicRepository)
	go notificationService.Run()

	return &App{
		Commands: Commands{
			ConversationService: conversationService,
			AuthService:         authService,
			NotificationService: notificationService,
			MessagingService:    messagingService,
		},
		Queries: Queries{
			Users:         usersRepository,
			Conversations: conversationsRepository,
			Messages:      messagesRepository,
		},
	}
}
