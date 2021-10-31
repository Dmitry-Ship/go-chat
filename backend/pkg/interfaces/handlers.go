package interfaces

import (
	"GitHub/go-chat/backend/pkg/application"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func HandleLogin(authService application.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokens, err := authService.Login(request.UserName, request.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			Expires:  time.Now().Add(application.AccessTokenExpiration),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(application.RefreshTokenExpiration),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleLogout(authService application.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		err := authService.Logout(userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleSignUp(authService application.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokens, err := authService.SignUp(request.UserName, request.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(application.AccessTokenExpiration),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(application.RefreshTokenExpiration),
			SameSite: http.SameSiteNoneMode,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleRefreshToken(authService application.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refresh_token")

		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newAccessToken, err := authService.RefreshAccessToken(refreshToken.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(application.AccessTokenExpiration),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleGetUser(authService application.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		user, err := authService.GetUser(userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)

	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := os.Getenv("ORIGIN_URL")

		return r.Header.Get("Origin") == origin
	},
}

func HandleWS(
	incomingMessageChannel chan<- json.RawMessage,
	registerClientChan chan<- *application.Client,
	unregisterClientChan chan<- *application.Client,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := application.NewClient(conn, unregisterClientChan, incomingMessageChannel, userID)

		registerClientChan <- client

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}

func HandleGetRooms(roomService application.RoomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := roomService.GetRooms()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(rooms)
	}
}

func HandleGetRoomsMessages(roomService application.RoomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomId, err := uuid.Parse(roomIdQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		messages, err := roomService.GetRoomMessages(roomId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		json.NewEncoder(w).Encode(data)
	}
}

func HandleGetRoom(roomService application.RoomService) http.HandlerFunc {
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

func HandleCreateRoom(roomService application.RoomService) http.HandlerFunc {
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

		err = roomService.CreateRoom(request.RoomId, request.RoomName, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleDeleteRoom(roomService application.RoomService) http.HandlerFunc {
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

		response := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{
			RoomId: request.RoomId,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func HandleJoinRoom(roomService application.RoomService) http.HandlerFunc {
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

		err = roomService.JoinRoom(request.RoomId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleLeaveRoom(roomService application.RoomService) http.HandlerFunc {
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

		err = roomService.LeaveRoom(request.RoomId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}
