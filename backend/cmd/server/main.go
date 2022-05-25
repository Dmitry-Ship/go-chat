package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/server"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
)

func createInfra() (context.Context, *app.Commands, readModel.QueriesRepository, infra.EventsSubscriber, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	redisClient := redisPubsub.GetRedisClient(ctx)
	db := postgres.NewDatabaseConnection()
	db.AutoMigrate()

	dbConnection := db.GetConnection()
	eventBus := infra.NewEventBus()

	activeClients := ws.NewActiveClients()

	commands := app.NewCommands(ctx, eventBus, dbConnection, activeClients, redisClient)
	queries := postgres.NewQueriesRepository(dbConnection)

	done := func() {
		cancel()
		redisClient.Close()
		eventBus.Close()
	}

	return ctx, commands, queries, eventBus, done
}

func main() {
	ctx, commands, queries, eventBus, done := createInfra()
	defer done()

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
