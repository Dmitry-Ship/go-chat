package commands

import (
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	"GitHub/go-chat/backend/internal/services"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewAuthCommands(ctx context.Context, eventPublisher infra.EventPublisher, db *gorm.DB) services.AuthService {
	usersRepository := postgres.NewUserRepository(db, eventPublisher)
	jwtTokens := services.NewJWTokens()

	return services.NewAuthService(usersRepository, jwtTokens)
}

func NewConversationCommands(ctx context.Context, eventPublisher infra.EventPublisher, db *gorm.DB) services.ConversationService {
	messagesRepository := postgres.NewMessageRepository(db, eventPublisher)
	groupConversationsRepository := postgres.NewGroupConversationRepository(db, eventPublisher)
	directConversationsRepository := postgres.NewDirectConversationRepository(db, eventPublisher)
	participantRepository := postgres.NewParticipantRepository(db, eventPublisher)

	return services.NewConversationService(groupConversationsRepository, directConversationsRepository, participantRepository, messagesRepository)
}

func NewNotificationsCommands(ctx context.Context, eventPublisher infra.EventPublisher, db *gorm.DB, redisClient *redis.Client) services.NotificationTopicService {
	notificationTopicRepository := postgres.NewNotificationTopicRepository(db, eventPublisher)

	return services.NewNotificationTopicService(ctx, notificationTopicRepository, redisClient)
}

func NewWSClientCommands(ctx context.Context, activeClients ws.ActiveClients, redisClient *redis.Client) services.ClientsService {
	return services.NewClientsService(ctx, redisClient, activeClients)
}
