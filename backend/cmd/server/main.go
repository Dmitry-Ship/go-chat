package main

import (
	"GitHub/go-chat/backend/internal/gracefulServer"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	"GitHub/go-chat/backend/internal/services"
	"context"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx, redisPubsub.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	defer redisClient.Close()

	db := postgres.NewDatabaseConnection(postgres.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})

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
	notificationService := services.NewNotificationService(ctx, redisClient)
	notificationResolverService := services.NewNotificationResolverService(participantRepository)
	queries := postgres.NewQueriesRepository(db)
	notificationBuilderService := services.NewNotificationBuilderService(queries)

	notificationPipelineService := services.NewNotificationsPipeline(notificationService, notificationResolverService, notificationBuilderService)

	server := server.NewServer(
		ctx,
		authService,
		conversationService,
		notificationPipelineService,
		notificationService,
		queries,
		eventBus,
	)
	server.Run()

	s := gracefulServer.NewGracefulServer()
	s.Run()
}
