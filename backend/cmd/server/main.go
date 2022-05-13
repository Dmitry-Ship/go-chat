package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	"GitHub/go-chat/backend/internal/services"
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

	commands := app.NewCommands(ctx, eventBus, dbConnection, activeClients)
	queries := postgres.NewQueriesRepository(dbConnection)

	handlers := httpHandlers.NewHTTPHandlers(commands, queries)
	handlers.InitRoutes()

	notificationTopicRepository := postgres.NewNotificationTopicRepository(dbConnection)
	notificationTopicService := services.NewNotificationTopicService(ctx, notificationTopicRepository, redisClient)

	messageEventHandlers := domainEventsHandlers.NewMessageEventHandlers(ctx, eventBus, commands)
	notificationsEventHandlers := domainEventsHandlers.NewNotificationsEventHandlers(ctx, eventBus, notificationTopicService, queries)

	go messageEventHandlers.ListenForEvents()
	go notificationsEventHandlers.ListenForEvents()

	broadcaster := ws.NewBroadcaster(ctx, redisClient, activeClients)
	go broadcaster.Run()

	server := server.NewGracefulServer()
	server.Run()
}
