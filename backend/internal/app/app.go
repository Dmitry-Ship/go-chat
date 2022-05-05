package app

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres"
	"GitHub/go-chat/backend/internal/services"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Commands struct {
	ConversationService      services.ConversationService
	AuthService              services.AuthService
	MessagingService         services.MessagingService
	NotificationTopicService services.NotificationTopicService
	ClientRegister           services.ClientRegister
}

func NewCommands(ctx context.Context, eventsPubSub domain.EventPublisher, redisClient *redis.Client, db *gorm.DB, activeClients ws.ActiveClients) *Commands {
	messagesRepository := postgres.NewMessageRepository(db, eventsPubSub)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db, eventsPubSub)
	participantRepository := postgres.NewParticipantRepository(db, eventsPubSub)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)

	return &Commands{
		ConversationService:      services.NewConversationService(conversationsRepository, participantRepository),
		AuthService:              services.NewAuthService(usersRepository),
		MessagingService:         services.NewMessagingService(messagesRepository),
		NotificationTopicService: services.NewNotificationTopicService(ctx, notificationTopicRepository, redisClient),
		ClientRegister:           services.NewClientRegister(activeClients),
	}
}
