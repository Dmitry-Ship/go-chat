package main

import (
	"GitHub/go-chat/backend/pkg/application"
	"GitHub/go-chat/backend/pkg/inmemory"
	"GitHub/go-chat/backend/pkg/interfaces"
	ws "GitHub/go-chat/backend/pkg/websocket"

	"log"
	"net/http"
	"os"
)

func main() {
	messagesRepository := inmemory.NewChatMessageRepository()
	usersRepository := inmemory.NewUserRepository()
	roomsRepository := inmemory.NewRoomRepository()
	participantRepository := inmemory.NewParticipantRepository()

	hub := ws.NewHub()

	authService := application.NewAuthService(usersRepository)
	roomCommandService := application.NewRoomCommandService(roomsRepository, participantRepository, usersRepository, messagesRepository, hub)
	roomQueryService := application.NewRoomQueryService(roomsRepository, participantRepository, usersRepository, messagesRepository)

	wsHandler := interfaces.NewWSMessageHandler(roomCommandService)

	ensureAuth := interfaces.MakeEnsureAuth(authService)

	http.HandleFunc("/ws", ensureAuth(interfaces.HandleWS(wsHandler.MessageChannel, hub.Register, hub.Unregister)))

	http.HandleFunc("/signup", interfaces.AddDefaultHeaders(interfaces.HandleSignUp(authService)))
	http.HandleFunc("/login", interfaces.AddDefaultHeaders(interfaces.HandleLogin(authService)))
	http.HandleFunc("/logout", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleLogout(authService))))
	http.HandleFunc("/refreshToken", interfaces.AddDefaultHeaders((interfaces.HandleRefreshToken(authService))))
	http.HandleFunc("/getUser", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetUser(authService))))
	http.HandleFunc("/getRooms", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetRooms(roomQueryService))))
	http.HandleFunc("/getRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetRoom(roomQueryService))))
	http.HandleFunc("/getRoomsMessages", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleGetRoomsMessages(roomQueryService))))
	http.HandleFunc("/createRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleCreateRoom(roomCommandService))))
	http.HandleFunc("/deleteRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleDeleteRoom(roomCommandService))))
	http.HandleFunc("/joinRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleJoinRoom(roomCommandService))))
	http.HandleFunc("/leaveRoom", interfaces.AddDefaultHeaders(ensureAuth(interfaces.HandleLeaveRoom(roomCommandService))))

	go wsHandler.Run()
	go hub.Run()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
