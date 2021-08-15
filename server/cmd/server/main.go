package main

import (
	"GitHub/go-chat/server/pkg/application"
	"GitHub/go-chat/server/pkg/interfaces"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	room := application.NewRoom()
	go room.Run()

	interfaces.HandleRequests(room)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
