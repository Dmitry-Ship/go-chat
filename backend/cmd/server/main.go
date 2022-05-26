package main

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/gracefulServer"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/infra/postgres"
	redisPubsub "GitHub/go-chat/backend/internal/infra/redis"
	"GitHub/go-chat/backend/internal/server"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redisPubsub.GetRedisClient(ctx)
	defer redisClient.Close()

	db := postgres.NewDatabaseConnection()
	db.AutoMigrate()

	dbConnection := db.GetConnection()
	eventBus := infra.NewEventBus()
	defer eventBus.Close()

	activeClients := ws.NewActiveClients()

	commands := app.NewCommands(ctx, eventBus, dbConnection, activeClients, redisClient)
	queries := postgres.NewQueriesRepository(dbConnection)

	server := server.NewServer(ctx, commands, queries, eventBus)

	server.InitRoutes()
	server.ListenForEvents()

	s := gracefulServer.NewGracefulServer()
	s.Run()
}
