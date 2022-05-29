package main

import (
	"GitHub/go-chat/backend/internal/commands"
	"GitHub/go-chat/backend/internal/gracefulServer"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	ws "GitHub/go-chat/backend/internal/websocket"
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

	activeClients := ws.NewActiveClients()

	authCommands := commands.NewAuthCommands(ctx, eventBus, dbConnection)
	conversationCommands := commands.NewConversationCommands(ctx, eventBus, dbConnection)
	notificationCommands := commands.NewNotificationsCommands(ctx, eventBus, dbConnection, redisClient)
	wsClientCommands := commands.NewWSClientCommands(ctx, activeClients, redisClient)
	queries := postgres.NewQueriesRepository(dbConnection)

	server := server.NewServer(ctx, authCommands, conversationCommands, notificationCommands, wsClientCommands, queries, eventBus)
	server.Run()

	s := gracefulServer.NewGracefulServer()
	s.Run()
}
