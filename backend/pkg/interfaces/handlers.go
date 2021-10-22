package interfaces

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/application"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func HandleRequests(
	roomService application.RoomService,
	authService application.AuthService,
	hub application.Hub,
	wsMessageChannel chan json.RawMessage,
) {
	http.HandleFunc("/signup", AddDefaultHeaders(handleSignUp(authService)))
	http.HandleFunc("/login", AddDefaultHeaders(handleLogin(authService)))
	http.HandleFunc("/logout", AddDefaultHeaders(EnsureAuth(handleLogout(authService), authService)))
	http.HandleFunc("/refreshToken", AddDefaultHeaders((handleRefreshToken(authService))))

	http.HandleFunc("/ws", EnsureAuth(handeleWS(wsMessageChannel, hub), authService))
	http.HandleFunc("/getRooms", AddDefaultHeaders(EnsureAuth(handleGetRooms(roomService), authService)))
	http.HandleFunc("/getRoom", AddDefaultHeaders(EnsureAuth(handleGetRoom(roomService), authService)))
	http.HandleFunc("/getUser", AddDefaultHeaders(EnsureAuth(handleGetUser(authService), authService)))
	http.HandleFunc("/getRoomsMessages", AddDefaultHeaders(EnsureAuth(handleGetRoomsMessages(roomService), authService)))
	http.HandleFunc("/createRoom", AddDefaultHeaders(EnsureAuth(handleCreateRoom(roomService), authService)))
	http.HandleFunc("/deleteRoom", AddDefaultHeaders(EnsureAuth(handleDeleteRoom(roomService), authService)))
}

func handleLogin(authService application.AuthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			UserName string `json:"user_name"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			log.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		tokens, err := authService.Login(request.UserName, request.Password)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			Expires:  time.Now().Add(application.AccessTokenExpiration),
			HttpOnly: true,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(application.RefreshTokenExpiration),
			HttpOnly: true,
			Secure:   true,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func handleLogout(authService application.AuthService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		err := authService.Logout(userID)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
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

func handleSignUp(authService application.AuthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			UserName string `json:"user_name"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			log.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		tokens, err := authService.SignUp(request.UserName, request.Password)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(application.AccessTokenExpiration),
			SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(application.RefreshTokenExpiration),
			SameSite: http.SameSiteNoneMode,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func handleRefreshToken(authService application.AuthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refresh_token")

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newAccessToken, err := authService.RefreshAccessToken(refreshToken.Value)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		url := os.Getenv("ORIGIN_URL")

		url = url[8:]

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(application.AccessTokenExpiration),
			SameSite: http.SameSiteNoneMode,
		})

		json.NewEncoder(w).Encode("OK")
	}
}

func handleGetUser(authService application.AuthService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		user, err := authService.GetUser(userID)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		json.NewEncoder(w).Encode(user)

	}
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
	wsMessageChannel chan json.RawMessage,
	hub application.Hub,
) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		client := application.NewClient(conn, hub, wsMessageChannel, userID)

		hub.Register(client)

		go client.SendNotifications()
		go client.ReceiveMessages()
	}
}

func handleGetRooms(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		rooms, err := roomService.GetRooms()

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		json.NewEncoder(w).Encode(rooms)
	}
}

func handleGetRoomsMessages(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
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

		json.NewEncoder(w).Encode(data)
	}
}

func handleGetRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		query := r.URL.Query()

		roomIdQuery := query.Get("room_id")
		roomId, err := uuid.Parse(roomIdQuery)

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
			Joined: roomService.HasJoined(roomId, userID),
		}

		json.NewEncoder(w).Encode(data)
	}
}

func handleCreateRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		request := struct {
			RoomName string    `json:"room_name"`
			RoomId   uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			log.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		err = roomService.CreateRoom(request.RoomId, request.RoomName, userID)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func handleDeleteRoom(roomService application.RoomService) func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	return func(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
		request := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			log.Println("Body parse error", err)
			w.WriteHeader(400)
			return
		}

		err = roomService.DeleteRoom(request.RoomId)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		json.NewEncoder(w).Encode("OK")

		response := struct {
			RoomId uuid.UUID `json:"room_id"`
		}{
			RoomId: request.RoomId,
		}

		json.NewEncoder(w).Encode(response)
	}
}
