package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
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
	err := db.AutoMigrate()

	if err != nil {
		panic(err)
	}

	dbConnection := db.GetConnection()
	domainEventsPubSub := domain.NewPubsub()

	commands := app.NewCommands(ctx, domainEventsPubSub, dbConnection)
	queries := postgres.NewQueriesRepository(dbConnection)

	notificationTopicRepository := postgres.NewNotificationTopicRepository(dbConnection)
	hub := ws.NewHub(ctx, redisClient, commands, notificationTopicRepository)
	go hub.Run()

	handlers := httpHandlers.NewHTTPHandlers(commands, queries, hub)
	handlers.InitRoutes()

	eventHandlers := domainEventsHandlers.NewEventHandlers(domainEventsPubSub, commands, queries, hub)
	eventHandlers.ListerForEvents()

	server := server.NewGracefulServer()
	server.Run()
}
