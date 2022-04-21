package main

import (
	"GitHub/go-chat/backend/pkg/app"
	"GitHub/go-chat/backend/pkg/httpServer"
	"GitHub/go-chat/backend/pkg/postgres"
	pubsub "GitHub/go-chat/backend/pkg/redis"
	ws "GitHub/go-chat/backend/pkg/websocket"
	"log"
	"net/http"
	"os"
)

func main() {
	redisClient := pubsub.GetRedisClient()
	connectionsHub := ws.NewHub(redisClient)
	go connectionsHub.Run()

	db := postgres.GetDatabaseConnection()
	application := app.NewApp(db, connectionsHub)
	server := httpServer.NewHTTPServer(application, connectionsHub)
	server.Init()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
