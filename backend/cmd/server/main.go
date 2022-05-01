package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
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
	ps := domain.NewPubsub()

	application := app.NewApp(ctx, ps, dbConnection, redisClient)

	queryController := httpHandlers.NewQueryController(&application.Queries)
	commandController := httpHandlers.NewCommandController(&application.Commands)
	handlers := httpHandlers.NewHTTPHandlers(queryController, commandController)
	handlers.InitRoutes()

	eventHandlers := domainEventsHandlers.NewEventHandlers(ps, &application.Commands, &application.Queries)
	eventHandlers.ListerForEvents()

	server := server.NewGracefulServer()
	server.Run()
}
