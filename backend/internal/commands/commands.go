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
	usersRepository := postgres.NewUserRepository(db, eventPublisher)

	return services.NewConversationService(groupConversationsRepository, directConversationsRepository, participantRepository, usersRepository, messagesRepository)
}

func NewNotificationsCommands(ctx context.Context, eventPublisher infra.EventPublisher, db *gorm.DB, redisClient *redis.Client) services.NotificationService {
	activeClients := ws.NewActiveClients()
	participantRepository := postgres.NewParticipantRepository(db, eventPublisher)

	return services.NewNotificationService(ctx, participantRepository, activeClients, redisClient)
}
