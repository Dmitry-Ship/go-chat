package main

import (
	"GitHub/go-chat/backend/pkg/application"
	"GitHub/go-chat/backend/pkg/database"
	"GitHub/go-chat/backend/pkg/interfaces"
	"GitHub/go-chat/backend/pkg/postgres"
	"time"

	ws "GitHub/go-chat/backend/pkg/websocket"

	"log"
	"net/http"
	"os"
)

func main() {
	db := database.GetDatabaseConnection()

	messagesRepository := postgres.NewChatMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	roomsRepository := postgres.NewRoomRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)

	const RefreshTokenExpiration = 24 * 90 * time.Hour
	const AccessTokenExpiration = 10 * time.Minute

	hub := ws.NewHub()

	roomCommandService := application.NewRoomCommandService(roomsRepository, participantRepository, usersRepository, messagesRepository, hub)
	roomQueryService := application.NewRoomQueryService(roomsRepository, participantRepository, usersRepository, messagesRepository)
	authService := application.NewAuthService(usersRepository, RefreshTokenExpiration, AccessTokenExpiration)
	ensureAuth := interfaces.MakeEnsureAuth(authService)

	hub.SetWSHandler("message", interfaces.HandleWSMessage(roomCommandService))

	http.HandleFunc("/signup", interfaces.AddDefaultHeaders(interfaces.HandleSignUp(authService)))
	http.HandleFunc("/login", interfaces.AddDefaultHeaders(interfaces.HandleLogin(authService)))
	http.HandleFunc("/logout", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleLogout(authService))))
	http.HandleFunc("/refreshToken", interfaces.AddDefaultHeaders((interfaces.HandleRefreshToken(authService))))
	http.HandleFunc("/getUser", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetUser(authService))))

	http.HandleFunc("/ws", ensureAuth(interfaces.HandleWS(hub)))
	http.HandleFunc("/getRooms", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetRooms(roomQueryService))))
	http.HandleFunc("/getRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetRoom(roomQueryService))))
	http.HandleFunc("/getRoomsMessages", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetRoomsMessages(roomQueryService))))
	http.HandleFunc("/createRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleCreateRoom(roomCommandService))))
	http.HandleFunc("/deleteRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleDeleteRoom(roomCommandService))))
	http.HandleFunc("/joinRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleJoinRoom(roomCommandService))))
	http.HandleFunc("/leaveRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleLeaveRoom(roomCommandService))))

	go hub.Run()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
