package main

import (
	"GitHub/go-chat/backend/pkg/application"
	"GitHub/go-chat/backend/pkg/inmemory"
	"GitHub/go-chat/backend/pkg/interfaces"
	"log"
	"net/http"
	"os"
)

func main() {
	messagesRepository := inmemory.NewChatMessageRepository()
	usersRepository := inmemory.NewUserRepository()
	roomsRepository := inmemory.NewRoomRepository()
	participantRepository := inmemory.NewParticipantRepository()

	hub := application.NewHub()
	go hub.Run()
	authService := application.NewAuthService(usersRepository)
	roomService := application.NewRoomService(roomsRepository, participantRepository, usersRepository, messagesRepository, hub)

	wsHandler := interfaces.NewWSMessageHandler(roomService)
	go wsHandler.Run()

	interfaces.HandleRequests(roomService, authService, hub, wsHandler.MessageChannel)

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
