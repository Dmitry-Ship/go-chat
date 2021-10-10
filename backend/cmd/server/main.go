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
)

func main() {
	messagesRepository := inmemory.NewChatMessageRepository()
	usersRepository := inmemory.NewUserRepository()
	roomsRepository := inmemory.NewRoomRepository()
	participantRepository := inmemory.NewParticipantRepository()

	userService := application.NewUserService(usersRepository)
	notificationService := application.NewNotificationService(participantRepository, userService)
	messageService := application.NewMessageService(messagesRepository, usersRepository, notificationService.Broadcast)
	roomService := application.NewRoomService(roomsRepository, participantRepository, userService, messageService)

	roomsRepository.Create(domain.NewRoom("Defalt Room"))

	go notificationService.Run()

	interfaces.HandleRequests(userService, messageService, roomService)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
