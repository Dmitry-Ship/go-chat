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
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func validateConfig() error {
	requiredVars := []string{
		"ACCESS_TOKEN_SECRET",
		"REFRESH_TOKEN_SECRET",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USER",
		"DB_PASSWORD",
		"REDIS_HOST",
		"REDIS_PORT",
		"CLIENT_ORIGIN",
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%s environment variable is required", envVar)
		}
	}

	if os.Getenv("ACCESS_TOKEN_SECRET") == "generate-with-make-secret-or-crypto-rand" {
		return fmt.Errorf("ACCESS_TOKEN_SECRET must be set to a strong secret (run 'make secret' to generate one)")
	}

	if os.Getenv("REFRESH_TOKEN_SECRET") == "generate-with-make-secret-or-crypto-rand" {
		return fmt.Errorf("REFRESH_TOKEN_SECRET must be set to a strong secret (run 'make secret' to generate one)")
	}

	if os.Getenv("REDIS_PASSWORD") == "change-this-redis-password" {
		log.Println("WARNING: REDIS_PASSWORD should be changed from the default value")
	}

	maxUserConnections, err := strconv.Atoi(os.Getenv("WS_RATE_LIMIT_MAX_USER"))
	if err != nil {
		return fmt.Errorf("WS_RATE_LIMIT_MAX_USER must be a valid integer")
	}
	if maxUserConnections <= 0 {
		return fmt.Errorf("WS_RATE_LIMIT_MAX_USER must be positive")
	}

	maxIPConnections, err := strconv.Atoi(os.Getenv("WS_RATE_LIMIT_MAX_IP"))
	if err != nil {
		return fmt.Errorf("WS_RATE_LIMIT_MAX_IP must be a valid integer")
	}
	if maxIPConnections <= 0 {
		return fmt.Errorf("WS_RATE_LIMIT_MAX_IP must be positive")
	}

	windowDurationStr := os.Getenv("WS_RATE_LIMIT_WINDOW")
	windowDuration, err := time.ParseDuration(windowDurationStr)
	if err != nil {
		return fmt.Errorf("WS_RATE_LIMIT_WINDOW must be a valid duration (e.g., '60s')")
	}
	if windowDuration <= 0 {
		return fmt.Errorf("WS_RATE_LIMIT_WINDOW must be positive")
	}

	return nil
}

func main() {
	if err := validateConfig(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

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

	pool, err := postgres.NewDatabaseConnection(ctx, postgres.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	cacheClient := cache.NewRedisCacheClient(redisClient, cache.CacheConfig{
		Prefix: "cache:",
	})

	messagesRepository := postgres.NewMessageRepository(pool, eventBus)
	groupConversationsRepository := postgres.NewGroupConversationRepository(pool, eventBus)
	directConversationsRepository := postgres.NewDirectConversationRepository(pool, eventBus)
	participantRepository := postgres.NewParticipantRepository(pool, eventBus)
	usersRepository := postgres.NewUserRepository(pool, eventBus)

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

	subscriptionSync := ws.NewSubscriptionSync(redisClient, redisPubsub.SubscriptionChannel)

	activeClients := ws.NewActiveClients()
	queries := postgres.NewQueriesRepository(pool)
	notificationService := services.NewNotificationServiceWithClients(ctx, redisClient, cachedParticipantRepository, subscriptionSync, queries, activeClients)

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
		notificationService,
		queries,
		eventBus,
		ipRateLimiter,
		userRateLimiter,
	)
	handler := server.Run()

	s := gracefulServer.NewGracefulServer(handler)
	s.Run()
}
