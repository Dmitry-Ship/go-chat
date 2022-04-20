package main

import (
	"GitHub/go-chat/backend/pkg/httpHandlers"
	"GitHub/go-chat/backend/pkg/postgres"
	pubsub "GitHub/go-chat/backend/pkg/redis"
	"GitHub/go-chat/backend/pkg/services"

	ws "GitHub/go-chat/backend/pkg/websocket"

	"log"
	"net/http"
	"os"
)

func main() {
	db := postgres.GetDatabaseConnection()
	redisClient := pubsub.GetRedisClient()
	hub := ws.NewHub(redisClient)

	messagesRepository := postgres.NewMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)

	conversationWSResolver := services.NewConversationWSResolver(participantRepository, messagesRepository, hub)

	conversationService := services.NewConversationService(conversationsRepository, participantRepository, messagesRepository, conversationWSResolver)
	authService := services.NewAuthService(usersRepository, usersRepository)

	ensureAuth := httpHandlers.MakeEnsureAuth(authService)

	wsHandlers := ws.NewWSHandlers()

	wsHandlers.SetWSHandler("message", httpHandlers.HandleWSMessage(conversationService))

	http.HandleFunc("/ws", ensureAuth(httpHandlers.HandleWS(hub, wsHandlers)))

	http.HandleFunc("/signup", httpHandlers.AddHeaders(httpHandlers.HandleSignUp(authService)))
	http.HandleFunc("/login", httpHandlers.AddHeaders(httpHandlers.HandleLogin(authService)))
	http.HandleFunc("/logout", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleLogout(authService))))
	http.HandleFunc("/refreshToken", httpHandlers.AddHeaders((httpHandlers.HandleRefreshToken(authService))))
	http.HandleFunc("/getUser", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleGetUser(usersRepository))))

	http.HandleFunc("/getConversations", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleGetConversations(conversationsRepository))))
	http.HandleFunc("/getContacts", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleGetContacts(usersRepository))))
	http.HandleFunc("/getConversation", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleGetConversation(conversationsRepository))))
	http.HandleFunc("/getConversationsMessages", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleGetConversationsMessages(messagesRepository))))

	http.HandleFunc("/createConversation", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleCreateConversation(conversationService))))
	http.HandleFunc("/deleteConversation", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleDeleteConversation(conversationService))))
	http.HandleFunc("/joinConversation", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleJoinPublicConversation(conversationService))))
	http.HandleFunc("/leaveConversation", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleLeavePublicConversation(conversationService))))
	http.HandleFunc("/renameConversation", httpHandlers.AddHeaders(ensureAuth(httpHandlers.HandleRenamePublicConversation(conversationService))))

	go hub.Run()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
