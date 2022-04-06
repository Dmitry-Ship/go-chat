package interfaces

import (
	"GitHub/go-chat/backend/pkg/application"
	ws "GitHub/go-chat/backend/pkg/websocket"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := os.Getenv("API_URL")

		return r.Header.Get("Origin") == origin
	},
}

func HandleWSMessage(roomService application.RoomCommandService) ws.WSHandler {
	return func(message ws.IncomingNotification, data json.RawMessage) {
		request := struct {
			Content string    `json:"content"`
			RoomId  uuid.UUID `json:"room_id"`
		}{}

		if err := json.Unmarshal([]byte(data), &request); err != nil {
			log.Println(err)
			return
		}

		err := roomService.SendMessage(request.Content, "user", request.RoomId, message.UserID)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func HandleWS(hub ws.Hub, wsHandlers ws.WSHandlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WS", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := ws.NewClient(conn, hub, wsHandlers, userID)

		hub.RegisterClient(client)

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}

func HandleGetRooms(roomService application.RoomQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := roomService.GetRooms()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(rooms)
	}
}

func HandleGetRoomsMessages(roomService application.RoomQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomId, err := uuid.Parse(roomIdQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, _ := r.Context().Value("userId").(uuid.UUID)

		messages, err := roomService.GetRoomMessages(roomId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Messages []application.MessageFull `json:"messages"`
		}{
			Messages: messages,
		}

		json.NewEncoder(w).Encode(data)
	}
}

func HandleGetRoom(roomService application.RoomQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomId, err := uuid.Parse(roomIdQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		room, err := roomService.GetRoom(roomId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(room)
	}
}

func HandleCreateRoom(roomService application.RoomCommandService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RoomName string    `json:"room_name"`
			RoomId   uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, _ := r.Context().Value("userId").(uuid.UUID)

		err = roomService.CreatePublicRoom(request.RoomId, request.RoomName, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleDeleteRoom(roomService application.RoomCommandService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = roomService.DeleteRoom(request.RoomId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleJoinPublicRoom(roomService application.RoomCommandService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, _ := r.Context().Value("userId").(uuid.UUID)

		err = roomService.JoinPublicRoom(request.RoomId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleLeavePublicRoom(roomService application.RoomCommandService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		request := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = roomService.LeavePublicRoom(request.RoomId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}