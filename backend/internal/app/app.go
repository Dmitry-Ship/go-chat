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
	ConversationService services.ConversationService
	AuthService         services.AuthService
	ClientsService      services.ClientsService
	NotificationService services.NotificationTopicService
}

func NewCommands(ctx context.Context, eventPublisher infra.EventPublisher, db *gorm.DB, activeClients ws.ActiveClients, redisClient *redis.Client) *Commands {
	messagesRepository := postgres.NewMessageRepository(db, eventPublisher)
	usersRepository := postgres.NewUserRepository(db)
	groupConversationsRepository := postgres.NewGroupConversationRepository(db, eventPublisher)
	directConversationsRepository := postgres.NewDirectConversationRepository(db, eventPublisher)
	participantRepository := postgres.NewParticipantRepository(db, eventPublisher)
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db)
	jwtTokens := services.NewJWTokens()

	return &Commands{
		ConversationService: services.NewConversationService(groupConversationsRepository, directConversationsRepository, participantRepository, messagesRepository),
		AuthService:         services.NewAuthService(usersRepository, jwtTokens),
		ClientsService:      services.NewClientsService(ctx, redisClient, activeClients),
		NotificationService: services.NewNotificationTopicService(ctx, notificationTopicRepository, redisClient),
	}
}
