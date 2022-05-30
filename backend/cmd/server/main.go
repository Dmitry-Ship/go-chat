package main

import (
	"GitHub/go-chat/backend/internal/commands"
	"GitHub/go-chat/backend/internal/gracefulServer"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx)
	defer redisClient.Close()

	db := postgres.NewDatabaseConnection()
	db.AutoMigrate()
	dbConnection := db.GetConnection()
	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	authCommands := commands.NewAuthCommands(ctx, eventBus, dbConnection)
	conversationCommands := commands.NewConversationCommands(ctx, eventBus, dbConnection)
	notificationCommands := commands.NewNotificationsCommands(ctx, eventBus, dbConnection, redisClient)
	queries := postgres.NewQueriesRepository(dbConnection)

	server := server.NewServer(ctx, authCommands, conversationCommands, notificationCommands, queries, eventBus)
	server.Run()

	s := gracefulServer.NewGracefulServer()
	s.Run()
}
