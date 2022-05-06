package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx)
	db := postgres.NewDatabaseConnection()
	db.AutoMigrate()

	dbConnection := db.GetConnection()
	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	activeClients := ws.NewActiveClients()

	commands := app.NewCommands(ctx, eventBus, redisClient, dbConnection, activeClients)
	queries := postgres.NewQueriesRepository(dbConnection)

	broadcaster := ws.NewBroadcaster(ctx, redisClient, activeClients)
	go broadcaster.Run()

	handlers := httpHandlers.NewHTTPHandlers(commands, queries)
	handlers.InitRoutes()

	eventHandlers := domainEventsHandlers.NewEventHandlers(ctx, eventBus, commands, queries)
	go eventHandlers.ListenForEvents()

	server := server.NewGracefulServer()
	server.Run()
}
