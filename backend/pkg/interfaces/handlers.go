package interfaces

import (
	"GitHub/go-chat/backend/common"
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/application"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func HandleRequests(
	userService application.UserService,
	messageService application.MessageService,
	roomService application.RoomService,
	hub application.Hub,
) {
	http.HandleFunc("/ws", handeleWS(userService, messageService, roomService, hub))
	http.HandleFunc("/getRooms", common.AddDefaultHeaders(handleGetRooms(roomService)))
	http.HandleFunc("/getRoomsMessages", common.AddDefaultHeaders(handleRoomsMessages(messageService, roomService)))
	http.HandleFunc("/createRoom", common.AddDefaultHeaders(handleCreateRoom(roomService)))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		clientURL := os.Getenv("ORIGIN_URL")

		return r.Header.Get("Origin") == clientURL
	},
}

func handeleWS(
	userService application.UserService,
	messageService application.MessageService,
	roomService application.RoomService,
	hub application.Hub,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user, err := userService.CreateUser(domain.NewUser())

		if err != nil {
			log.Println(err)
			return
		}

		client := application.NewClient(conn, messageService, userService, roomService, user.Id)

		hub.Register(client)

		data := struct {
			UserId int32 `json:"user_id"`
		}{
			UserId: user.Id,
		}

		client.Send <- hub.NewNotification("user_id", data)

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}

func handleGetRooms(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := roomService.GetRooms()

		if err != nil {
			log.Println(err)
			return
		}

		common.SendJSONresponse(rooms, w)
	}
}

func handleRoomsMessages(messageService application.MessageService, roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomIdInt, err := strconv.ParseInt(roomIdQuery, 0, 32)

		if err != nil {
			log.Println(err)
			return
		}
		roomIdInt32 := int32(roomIdInt)

		userId := query.Get("user_id")
		userIdInt, err := strconv.ParseInt(userId, 0, 32)

		if err != nil {
			log.Println(err)
			return
		}
		userIdInt32 := int32(userIdInt)

		room, err := roomService.GetRoom(roomIdInt32)

		if err != nil {
			log.Println(err)
			return
		}

		messages, err := messageService.GetMessagesFull(roomIdInt32)

		if err != nil {
			log.Println(err)
			return
		}

		messagesValue := []application.MessageFull{}
		for _, message := range messages {
			messagesValue = append(messagesValue, *message)
		}

		data := struct {
			Room     domain.Room               `json:"room"`
			Joined   bool                      `json:"joined"`
			Messages []application.MessageFull `json:"messages"`
		}{
			Room:     *room,
			Joined:   roomService.HasJoined(userIdInt32, roomIdInt32),
			Messages: messagesValue,
		}

		common.SendJSONresponse(data, w)
	}
}

func handleCreateRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RoomName string `json:"room_name"`
			UserId   int32  `json:"user_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			fmt.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		room, err := roomService.CreateRoom(request.RoomName, request.UserId)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		common.SendJSONresponse(room, w)
	}
}
