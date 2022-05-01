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
	redisClient := redisPubsub.GetRedisClient()
	db := postgres.NewDatabaseConnection()
	err := db.AutoMigrate()

	if err != nil {
		panic(err)
	}

	dbConnection := db.GetConnection()
	ctx := context.Background()
	domainEventsPubSub := domain.NewPubsub()
	wsConnectionsPool := ws.NewConnectionsPool(ctx, redisClient)

	app := app.NewApp(ctx, domainEventsPubSub, dbConnection, wsConnectionsPool)

	wsHandlers := httpHandlers.NewWSHandlers(&app.Commands, wsConnectionsPool.IncomingNotifications)
	handlers := httpHandlers.NewHTTPHandlers(app)
	eventHandlers := domainEventsHandlers.NewEventHandlers(domainEventsPubSub, app)

	eventHandlers.ListerForEvents()
	handlers.InitRoutes()
	go wsHandlers.Run()
	go wsConnectionsPool.Run()

	server := server.NewGracefulServer()
	server.Run()
}
