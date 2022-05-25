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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx)
	db := postgres.NewDatabaseConnection()
	db.AutoMigrate()

	dbConnection := db.GetConnection()
	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	activeClients := ws.NewActiveClients()

	commands := app.NewCommands(ctx, eventBus, dbConnection, activeClients, redisClient)
	queries := postgres.NewQueriesRepository(dbConnection)

	handlers := httpHandlers.NewHTTPHandlers(commands, queries)
	handlers.InitRoutes()

	messageEventHandlers := domainEventsHandlers.NewMessageEventHandlers(ctx, eventBus, commands)
	notificationsEventHandlers := domainEventsHandlers.NewNotificationsEventHandlers(ctx, eventBus, commands, queries)

	go messageEventHandlers.ListenForEvents()
	go notificationsEventHandlers.ListenForEvents()
	go commands.ClientsService.Run()

	server := server.NewGracefulServer()
	server.Run()
}
