package main

import (
	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/gracefulServer"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/cache"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/ratelimit"
	"GitHub/go-chat/backend/internal/server"
	"GitHub/go-chat/backend/internal/services"
	"context"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx, redisPubsub.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	defer func() {
		_ = redisClient.Close()
	}()

	db := postgres.NewDatabaseConnection(postgres.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	cacheClient := cache.NewRedisCacheClient(redisClient, cache.CacheConfig{
		Prefix: "cache:",
	})

	messagesRepository := postgres.NewMessageRepository(db, eventBus)
	groupConversationsRepository := postgres.NewGroupConversationRepository(db, eventBus)
	directConversationsRepository := postgres.NewDirectConversationRepository(db, eventBus)
	participantRepository := postgres.NewParticipantRepository(db, eventBus)
	usersRepository := postgres.NewUserRepository(db, eventBus)

	cachedUsersRepository := cache.NewUserCacheDecorator(usersRepository, cacheClient)
	cachedGroupConversationsRepository := cache.NewGroupConversationCacheDecorator(groupConversationsRepository, cacheClient)
	cachedParticipantRepository := cache.NewParticipantCacheDecorator(participantRepository, cacheClient)

	authService := services.NewAuthService(cachedUsersRepository, config.Auth{
		AccessToken:  config.Token{Secret: os.Getenv("ACCESS_TOKEN_SECRET"), TTL: 10 * time.Minute},
		RefreshToken: config.Token{Secret: os.Getenv("REFRESH_TOKEN_SECRET"), TTL: 24 * 90 * time.Hour},
	})

	conversationService := services.NewConversationService(
		cachedGroupConversationsRepository,
		directConversationsRepository,
		cachedParticipantRepository,
		cachedUsersRepository,
		messagesRepository,
	)
	notificationService := services.NewNotificationService(ctx, redisClient)
	notificationResolverService := services.NewNotificationResolverService(cachedParticipantRepository)
	queries := postgres.NewQueriesRepository(db)
	notificationBuilderService := services.NewNotificationBuilderService(queries)

	notificationPipelineService := services.NewNotificationsPipeline(notificationService, notificationResolverService, notificationBuilderService)

	cacheInvalidationService := cache.NewCacheInvalidationService(cacheClient, eventBus)
	go cacheInvalidationService.Run(ctx)

	maxUserConnections, _ := strconv.Atoi(os.Getenv("WS_RATE_LIMIT_MAX_USER"))
	maxIPConnections, _ := strconv.Atoi(os.Getenv("WS_RATE_LIMIT_MAX_IP"))
	windowDurationStr := os.Getenv("WS_RATE_LIMIT_WINDOW")
	windowDuration, _ := time.ParseDuration(windowDurationStr)

	if maxUserConnections == 0 {
		maxUserConnections = 10
	}
	if maxIPConnections == 0 {
		maxIPConnections = 20
	}
	if windowDuration == 0 {
		windowDuration = 60 * time.Second
	}

	userRateLimiter := ratelimit.NewSlidingWindowRateLimiter(ratelimit.Config{
		MaxConnections: maxUserConnections,
		WindowDuration: windowDuration,
	})

	ipRateLimiter := ratelimit.NewSlidingWindowRateLimiter(ratelimit.Config{
		MaxConnections: maxIPConnections,
		WindowDuration: windowDuration,
	})

	server := server.NewServer(
		ctx,
		config.ServerConfig{
			ClientOrigin: os.Getenv("CLIENT_ORIGIN"),
		},
		authService,
		conversationService,
		notificationPipelineService,
		notificationService,
		queries,
		eventBus,
		ipRateLimiter,
		userRateLimiter,
	)
	server.Run()

	s := gracefulServer.NewGracefulServer()
	s.Run()
}
