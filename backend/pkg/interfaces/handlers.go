package interfaces

import (
	"GitHub/go-chat/backend/pkg/application"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleRequests(room *application.Room) {
	http.HandleFunc("/ws", handeleWS(room))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handeleWS(room *application.Room) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := application.NewClient(conn, room)

		data := struct {
			ClientId string `json:"client_id"`
		}{
			ClientId: client.Id,
		}

		client.Send <- client.NewNotification("client_id", data)

		client.Room.Join <- client

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}
