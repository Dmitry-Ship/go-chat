package main

import (
	"GitHub/go-chat/backend/internal/gracefulServer"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	"GitHub/go-chat/backend/internal/services"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx)
	defer redisClient.Close()

	db := postgres.NewDatabaseConnection()
	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	messagesRepository := postgres.NewMessageRepository(db, eventBus)
	groupConversationsRepository := postgres.NewGroupConversationRepository(db, eventBus)
	directConversationsRepository := postgres.NewDirectConversationRepository(db, eventBus)
	participantRepository := postgres.NewParticipantRepository(db, eventBus)
	usersRepository := postgres.NewUserRepository(db, eventBus)

	authService := services.NewAuthService(usersRepository)
	conversationService := services.NewConversationService(
		groupConversationsRepository,
		directConversationsRepository,
		participantRepository,
		usersRepository,
		messagesRepository,
	)
	notificationService := services.NewNotificationService(ctx, participantRepository, redisClient)
	queries := postgres.NewQueriesRepository(db)

	server := server.NewServer(ctx, authService, conversationService, notificationService, queries, eventBus)
	server.Run()

	s := gracefulServer.NewGracefulServer()
	s.Run()
}
