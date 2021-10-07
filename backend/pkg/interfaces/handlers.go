package interfaces

import (
	"GitHub/go-chat/backend/common"
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/application"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleRequests(hub *domain.Hub) {
	http.HandleFunc("/ws", handeleWS(hub))
	http.HandleFunc("/getRooms", handeleGetRooms)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.URL.Path == "/ws"
	},
}

func handeleWS(hub *domain.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := domain.NewUser()

		client := application.NewClient(conn, hub, user)

		data := struct {
			UserId int64 `json:"user_id"`
		}{
			UserId: user.Id,
		}

		client.User.Send <- user.NewNotification("user_id", data)
		client.Hub.Join <- domain.NewSubscription(user, 123)

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}

func handeleGetRooms(w http.ResponseWriter, r *http.Request) {

	defaultRoom := domain.NewRoom("default")
	rooms := []domain.Room{*defaultRoom}

	common.SendJSONresponse(rooms, w)
}
