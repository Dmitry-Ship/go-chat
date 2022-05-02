package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"GitHub/go-chat/backend/internal/server"
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
	wsConnectionsPool := ws.NewConnectionsPool(ctx, redisClient)
	go wsConnectionsPool.Run()

	app := app.NewApp(ctx, domainEventsPubSub, dbConnection, wsConnectionsPool)

	wsHandlers := httpHandlers.NewWSHandlers(&app.Commands, wsConnectionsPool.IncomingNotifications)
	go wsHandlers.Run()

	handlers := httpHandlers.NewHTTPHandlers(app)
	handlers.InitRoutes()

	eventHandlers := domainEventsHandlers.NewEventHandlers(domainEventsPubSub, app)
	eventHandlers.ListerForEvents()

	server := server.NewGracefulServer()
	server.Run()
}
