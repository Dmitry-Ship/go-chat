package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/httpServer"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"log"
	"net/http"
	"os"
)

func main() {
	redisClient := redisPubsub.GetRedisClient()
	websocketConnectionsHub := ws.NewHub(redisClient)
	go websocketConnectionsHub.Run()
	db := postgres.NewDatabaseConnection()
	db.RunMigrations()
	dbConnection := db.GetConnection()

	application := app.NewApp(dbConnection, websocketConnectionsHub)
	queryController := httpServer.NewQueryController(&application.Queries)
	commandController := httpServer.NewCommandController(&application.Commands)
	server := httpServer.NewHTTPServer(queryController, commandController)
	server.InitRoutes()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
