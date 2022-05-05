package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/hub"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
)

func main() {
	ctx := context.Background()
	redisClient := redisPubsub.GetRedisClient(ctx)
	db := postgres.NewDatabaseConnection()
	db.AutoMigrate()

	dbConnection := db.GetConnection()
	domainEventsPubSub := domain.NewPubsub()

	commands := app.NewCommands(ctx, domainEventsPubSub, redisClient, dbConnection)
	queries := postgres.NewQueriesRepository(dbConnection)

	activeClients := ws.NewActiveClients()

	broadcaster := hub.NewBroadcaster(ctx, redisClient, activeClients)
	go broadcaster.Run()
	clientRegister := hub.NewClientRegister(commands, activeClients)

	handlers := httpHandlers.NewHTTPHandlers(commands, queries, clientRegister)
	handlers.InitRoutes()

	eventHandlers := domainEventsHandlers.NewEventHandlers(domainEventsPubSub, commands, queries)
	eventHandlers.ListerForEvents()

	server := server.NewGracefulServer()
	server.Run()
}
