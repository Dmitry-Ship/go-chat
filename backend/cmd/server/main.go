package main

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/application"
	"GitHub/go-chat/backend/pkg/inmemory"
	"GitHub/go-chat/backend/pkg/interfaces"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	hub := application.NewHub()
	go hub.Run()

	messagesRepository := inmemory.NewChatMessageRepository()
	usersRepository := inmemory.NewUserRepository()
	roomsRepository := inmemory.NewRoomRepository()
	roomsRepository.Store(domain.NewRoom(uuid.New(), "Default Room"))
	participantRepository := inmemory.NewParticipantRepository()

	roomService := application.NewRoomService(roomsRepository, participantRepository, usersRepository, messagesRepository, hub)
	userService := application.NewUserService(usersRepository)

	wsHandler := interfaces.NewWSMessageHandler(userService, roomService)
	go wsHandler.Run()

	interfaces.HandleRequests(userService, roomService, hub, wsHandler.MessageChannel)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
