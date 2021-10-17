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

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func HandleRequests(
	userService application.UserService,
	roomService application.RoomService,
	hub application.Hub,
	wsMessageChannel chan json.RawMessage,
) {
	http.HandleFunc("/ws", handeleWS(userService, wsMessageChannel, hub))
	http.HandleFunc("/getRooms", common.AddDefaultHeaders(handleGetRooms(roomService)))
	http.HandleFunc("/getRoom", common.AddDefaultHeaders(handleGetRoom(roomService)))
	http.HandleFunc("/getRoomsMessages", common.AddDefaultHeaders(handleGetRoomsMessages(roomService)))
	http.HandleFunc("/createRoom", common.AddDefaultHeaders(handleCreateRoom(roomService)))
	http.HandleFunc("/deleteRoom", common.AddDefaultHeaders(handleDeleteRoom(roomService)))
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
	wsMessageChannel chan json.RawMessage,
	hub application.Hub,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := domain.NewUser()
		err = userService.CreateUser(user)

		if err != nil {
			log.Println(err)
			return
		}

		client := application.NewClient(conn, hub, wsMessageChannel, user.Id)

		hub.Register(client)

		data := struct {
			UserId uuid.UUID `json:"user_id"`
		}{
			UserId: user.Id,
		}

		client.SendNotification("user_id", data)

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

func handleGetRoomsMessages(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomId, err := uuid.Parse(roomIdQuery)

		if err != nil {
			log.Println(err)
			return
		}

		messages, err := roomService.GetRoomMessages(roomId)

		if err != nil {
			log.Println(err)
			return
		}

		messagesValue := []application.MessageFull{}
		for _, message := range messages {
			messagesValue = append(messagesValue, *message)
		}

		data := struct {
			Messages []application.MessageFull `json:"messages"`
		}{
			Messages: messagesValue,
		}

		common.SendJSONresponse(data, w)
	}
}

func handleGetRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomId, err := uuid.Parse(roomIdQuery)

		if err != nil {
			log.Println(err)
			return
		}

		userIdQuery := query.Get("user_id")
		userId, err := uuid.Parse(userIdQuery)

		if err != nil {
			log.Println(err)
			return
		}

		room, err := roomService.GetRoom(roomId)

		if err != nil {
			log.Println(err)
			return
		}

		data := struct {
			Room   domain.Room `json:"room"`
			Joined bool        `json:"joined"`
		}{
			Room:   *room,
			Joined: roomService.HasJoined(userId, roomId),
		}

		common.SendJSONresponse(data, w)
	}
}

func handleCreateRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RoomName string    `json:"room_name"`
			UserId   uuid.UUID `json:"user_id"`
			RoomId   uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			fmt.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		err = roomService.CreateRoom(request.RoomId, request.RoomName, request.UserId)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		common.SendJSONresponse(nil, w)
	}
}

func handleDeleteRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			fmt.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		err = roomService.DeleteRoom(request.RoomId)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)

		response := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{
			RoomId: request.RoomId,
		}

		common.SendJSONresponse(response, w)

	}
}
