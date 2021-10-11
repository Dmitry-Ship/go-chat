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
	roomsRepository.Create(domain.NewRoom("Default Room"))
	participantRepository := inmemory.NewParticipantRepository()

	notificationService := application.NewNotificationService(participantRepository)
	go notificationService.Run()

	userService := application.NewUserService(usersRepository)
	messageService := application.NewMessageService(messagesRepository, usersRepository, notificationService.Broadcast)
	roomService := application.NewRoomService(roomsRepository, participantRepository, usersRepository, messageService)

	interfaces.HandleRequests(userService, messageService, roomService, notificationService)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
