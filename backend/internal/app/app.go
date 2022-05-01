package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	"context"

	"gorm.io/gorm"
)

type App struct {
	Commands Commands
	Queries  readModel.QueriesRepository
}

type Commands struct {
	ConversationService services.ConversationService
	AuthService         services.AuthService
	NotificationService services.NotificationService
	MessagingService    services.MessagingService
}

func NewApp(ctx context.Context, eventsPubSub domain.EventPublisher, db *gorm.DB, connectionsPool ws.ConnectionsPool) *App {
	messagesRepository := postgres.NewMessageRepository(db, eventsPubSub)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db, eventsPubSub)
	participantRepository := postgres.NewParticipantRepository(db, eventsPubSub)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	return &App{
		Commands: Commands{
			ConversationService: services.NewConversationService(conversationsRepository, participantRepository),
			AuthService:         services.NewAuthService(usersRepository),
			NotificationService: services.NewNotificationService(connectionsPool, notificationTopicRepository),
			MessagingService:    services.NewMessagingService(messagesRepository),
		},
		Queries: postgres.NewQueriesRepository(db),
	}
}
