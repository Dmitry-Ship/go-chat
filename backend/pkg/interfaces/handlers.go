package interfaces

import (
	"GitHub/go-chat/backend/common"
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/application"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func HandleRequests(userService application.UserService, messageService application.MessageService, roomService application.RoomService) {
	http.HandleFunc("/ws", handeleWS(userService, messageService, roomService))
	http.HandleFunc("/getRooms", handleGetRooms(roomService))
	http.HandleFunc("/getRoomsMessages", handleRoomsMessages(messageService, roomService))
	http.HandleFunc("/createRoom", handleCreateRoom(roomService))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.URL.Path == "/ws"
	},
}

func handeleWS(userService application.UserService, messageService application.MessageService, roomService application.RoomService) func(w http.ResponseWriter, r *http.Request) {
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

		sendChan := userService.AddWSClient(user.Id)

		client := application.NewClient(conn, messageService, userService, roomService, sendChan)

		data := struct {
			UserId int32 `json:"user_id"`
		}{
			UserId: user.Id,
		}

		client.Send <- userService.NewNotification("user_id", data)

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
		roomId := query.Get("room_id")

		roomIdInt, err := strconv.ParseInt(roomId, 0, 32)

		if err != nil {
			log.Println(err)
			return
		}

		result := int32(roomIdInt)

		messages, err := messageService.GetMessagesFull(result)

		if err != nil {
			log.Println(err)
			return
		}

		room, err := roomService.GetRoom(result)

		if err != nil {
			log.Println(err)
			return
		}

		var messagesValue []application.MessageFull
		for _, message := range messages {
			messagesValue = append(messagesValue, *message)
		}

		data := struct {
			Room     domain.Room               `json:"room"`
			Messages []application.MessageFull `json:"messages"`
		}{
			Room:     *room,
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

		err := common.DecodeJSONBody(w, r, &request)
		if err != nil {
			var mr *common.MalformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.Msg, mr.Status)
			} else {
				log.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		room, err := roomService.CreateRoom(request.RoomName, request.UserId)

		if err != nil {
			log.Println(err)
			return
		}

		common.SendJSONresponse(room, w)
	}
}
