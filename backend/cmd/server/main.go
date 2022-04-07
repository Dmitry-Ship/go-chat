package main

import (
	"GitHub/go-chat/backend/pkg/application"
	"GitHub/go-chat/backend/pkg/database"
	"GitHub/go-chat/backend/pkg/interfaces"
	"GitHub/go-chat/backend/pkg/postgres"
	pubsub "GitHub/go-chat/backend/pkg/redis"

	ws "GitHub/go-chat/backend/pkg/websocket"

	"log"
	"net/http"
	"os"
)

func main() {
	db := database.GetDatabaseConnection()

	messagesRepository := postgres.NewChatMessageRepository(db)
	usersRepository := postgres.NewUserRepository(db)
	conversationsRepository := postgres.NewConversationRepository(db)
	participantRepository := postgres.NewParticipantRepository(db)

	redisClient := pubsub.GetRedisClient()
	hub := ws.NewHub(redisClient)

	conversationCommandService := application.NewConversationCommandService(conversationsRepository, participantRepository, usersRepository, messagesRepository, hub)
	conversationQueryService := application.NewConversationQueryService(conversationsRepository, participantRepository, usersRepository, messagesRepository)
	authService := application.NewAuthService(usersRepository)
	contactsQueryService := application.NewContactsQueryService(usersRepository)
	ensureAuth := interfaces.MakeEnsureAuth(authService)

	wsHandlers := ws.NewWSHandlers()
	wsHandlers.SetWSHandler("message", interfaces.HandleWSMessage(conversationCommandService))
	http.HandleFunc("/ws", ensureAuth(interfaces.HandleWS(hub, wsHandlers)))

	http.HandleFunc("/signup", interfaces.AddHeaders(interfaces.HandleSignUp(authService)))
	http.HandleFunc("/login", interfaces.AddHeaders(interfaces.HandleLogin(authService)))
	http.HandleFunc("/logout", interfaces.AddHeaders(ensureAuth(interfaces.HandleLogout(authService))))
	http.HandleFunc("/refreshToken", interfaces.AddHeaders((interfaces.HandleRefreshToken(authService))))
	http.HandleFunc("/getUser", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetUser(authService))))

	http.HandleFunc("/getConversations", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetConversations(conversationQueryService))))
	http.HandleFunc("/getContacts", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetContacts(contactsQueryService))))
	http.HandleFunc("/getConversation", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetConversation(conversationQueryService))))
	http.HandleFunc("/getConversationsMessages", interfaces.AddHeaders(ensureAuth(interfaces.HandleGetConversationsMessages(conversationQueryService))))
	http.HandleFunc("/createConversation", interfaces.AddHeaders(ensureAuth(interfaces.HandleCreateConversation(conversationCommandService))))
	http.HandleFunc("/deleteConversation", interfaces.AddHeaders(ensureAuth(interfaces.HandleDeleteConversation(conversationCommandService))))
	http.HandleFunc("/joinConversation", interfaces.AddHeaders(ensureAuth(interfaces.HandleJoinPublicConversation(conversationCommandService))))
	http.HandleFunc("/leaveConversation", interfaces.AddHeaders(ensureAuth(interfaces.HandleLeavePublicConversation(conversationCommandService))))

	go hub.Run()

	port := os.Getenv("PORT")

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
