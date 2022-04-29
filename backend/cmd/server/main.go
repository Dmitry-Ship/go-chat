package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
)

func main() {
	redisClient := redisPubsub.GetRedisClient()
	db := postgres.NewDatabaseConnection()
	err := db.AutoMigrate()

	if err != nil {
		panic(err)
	}

	dbConnection := db.GetConnection()

	application := app.NewApp(dbConnection, redisClient)

	queryController := httpHandlers.NewQueryController(&application.Queries)
	commandController := httpHandlers.NewCommandController(&application.Commands)
	handlers := httpHandlers.NewHTTPHandlers(queryController, commandController)
	handlers.InitRoutes()

	server := server.NewGracefulServer()
	server.Run()
}
