package main

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/interfaces"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	hub := domain.NewHub()
	go hub.Run()

	interfaces.HandleRequests(hub)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
