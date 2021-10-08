package interfaces

import (
	"GitHub/go-chat/backend/common"
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/application"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func HandleRequests(hub *domain.Hub) {
	http.HandleFunc("/ws", handeleWS(hub))
	http.HandleFunc("/getRooms", handleGetRooms(hub))
	http.HandleFunc("/getRoomsMessages", handleRoomsMessages(hub))
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
			UserId int32 `json:"user_id"`
		}{
			UserId: user.Id,
		}

		client.User.Send <- user.NewNotification("user_id", data)

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}

func handleGetRooms(hub *domain.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms := make([]domain.Room, 0, len(hub.Rooms))

		for _, value := range hub.Rooms {
			rooms = append(rooms, *value)
		}

		common.SendJSONresponse(rooms, w)
	}
}

func handleRoomsMessages(hub *domain.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		roomId := query.Get("room_id")

		roomIdInt, err := strconv.ParseInt(roomId, 0, 32)

		if err != nil {
			log.Println(err)
			return
		}

		result := int32(roomIdInt)

		data := struct {
			Room     domain.Room          `json:"room"`
			Messages []domain.ChatMessage `json:"messages"`
		}{
			Room:     *hub.Rooms[result],
			Messages: []domain.ChatMessage{},
		}

		common.SendJSONresponse(data, w)
	}
}
