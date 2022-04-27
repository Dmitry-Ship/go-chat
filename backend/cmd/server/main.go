package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/httpServer"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"log"
	"net/http"
	"os"
)

func main() {
	redisClient := redisPubsub.GetRedisClient()
	db := postgres.NewDatabaseConnection()
	db.RunMigrations()
	dbConnection := db.GetConnection()

	application := app.NewApp(dbConnection, redisClient)
	queryController := httpServer.NewQueryController(&application.Queries)
	commandController := httpServer.NewCommandController(&application.Commands)
	server := httpServer.NewHTTPServer(queryController, commandController)
	server.InitRoutes()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
