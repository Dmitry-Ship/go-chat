package interfaces

import (
	"GitHub/go-chat/pkg/application"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleRequests(hub *application.Hub) {
	http.Handle("/", http.FileServer(http.Dir("frontend")))
	http.HandleFunc("/ws", handeleWS(hub))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handeleWS(hub *application.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &application.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
		client.Hub.Register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WriteMessages()
		go client.Receive()
	}
}
