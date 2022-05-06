package app

import (
	"GitHub/go-chat/backend/internal/infra"
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

func NewCommands(ctx context.Context, eventPublisher infra.EventPublisher, redisClient *redis.Client, db *gorm.DB, activeClients ws.ActiveClients) *Commands {
	messagesRepository := postgres.NewMessageRepository(db, eventPublisher)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db, eventPublisher)
	participantRepository := postgres.NewParticipantRepository(db, eventPublisher)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)
	jwtTokens := services.NewJWTokens()

	return &Commands{
		ConversationService:      services.NewConversationService(conversationsRepository, participantRepository),
		AuthService:              services.NewAuthService(usersRepository, jwtTokens),
		MessagingService:         services.NewMessagingService(messagesRepository),
		NotificationTopicService: services.NewNotificationTopicService(ctx, notificationTopicRepository, redisClient),
		ClientRegister:           services.NewClientRegister(activeClients),
	}
}
