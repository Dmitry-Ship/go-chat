package main

import (
	"GitHub/go-chat/backend/pkg/application"
	"GitHub/go-chat/backend/pkg/database"
	"GitHub/go-chat/backend/pkg/interfaces"
	"GitHub/go-chat/backend/pkg/postgres"

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

	hub := ws.NewHub()
	wsHandlers := ws.NewWSHandlers()

	roomCommandService := application.NewRoomCommandService(roomsRepository, participantRepository, usersRepository, messagesRepository, hub)
	roomQueryService := application.NewRoomQueryService(roomsRepository, participantRepository, usersRepository, messagesRepository)
	authService := application.NewAuthService(usersRepository)
	contactsQueryService := application.NewContactsQueryService(usersRepository)
	ensureAuth := interfaces.MakeEnsureAuth(authService)

	wsHandlers.SetWSHandler("message", interfaces.HandleWSMessage(roomCommandService))

	http.HandleFunc("/signup", interfaces.AddHeaders(interfaces.HandleSignUp(authService)))
	http.HandleFunc("/login", interfaces.AddHeaders(interfaces.HandleLogin(authService)))
	http.HandleFunc("/logout", interfaces.AddHeaders(ensureAuth(interfaces.HandleLogout(authService))))
	http.HandleFunc("/refreshToken", interfaces.AddHeaders((interfaces.HandleRefreshToken(authService))))
	http.HandleFunc("/getUser", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetUser(authService))))

	http.HandleFunc("/ws", ensureAuth(interfaces.HandleWS(hub, wsHandlers)))

	http.HandleFunc("/getRooms", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetRooms(roomQueryService))))
	http.HandleFunc("/getContacts", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetContacts(contactsQueryService))))
	http.HandleFunc("/getRoom", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetRoom(roomQueryService))))
	http.HandleFunc("/getRoomsMessages", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetRoomsMessages(roomQueryService))))
	http.HandleFunc("/createRoom", interfaces.AddHeaders(ensureAuth(interfaces.HandleCreateRoom(roomCommandService))))
	http.HandleFunc("/deleteRoom", interfaces.AddHeaders(ensureAuth(interfaces.HandleDeleteRoom(roomCommandService))))
	http.HandleFunc("/joinRoom", interfaces.AddHeaders(ensureAuth(interfaces.HandleJoinRoom(roomCommandService))))
	http.HandleFunc("/leaveRoom", interfaces.AddHeaders(ensureAuth(interfaces.HandleLeaveRoom(roomCommandService))))

	go wsHandlers.Run()
	go hub.Run()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
