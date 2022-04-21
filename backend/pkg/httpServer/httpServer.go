package httpServer

import (
	"GitHub/go-chat/backend/pkg/app"
	ws "GitHub/go-chat/backend/pkg/websocket"
	"net/http"
)

type HTTPServer struct {
	app              app.App
	WSConnectionsHub ws.Hub
}

func NewHTTPServer(app *app.App, hub ws.Hub) *HTTPServer {
	return &HTTPServer{
		app:              *app,
		WSConnectionsHub: hub,
	}
}

func (s *HTTPServer) Init() {
	ensureAuth := MakeEnsureAuth(s.app.Commands.AuthService)

	wsHandlers := ws.NewWSHandlers()

	wsHandlers.SetWSHandler("message", s.handleWSMessage)

	http.HandleFunc("/ws", ensureAuth(s.handleWS(wsHandlers)))

	http.HandleFunc("/signup", addHeaders(s.handleSignUp))
	http.HandleFunc("/login", addHeaders(s.handleLogin))
	http.HandleFunc("/logout", addHeaders(ensureAuth(s.handleLogout)))
	http.HandleFunc("/refreshToken", addHeaders((s.handleRefreshToken)))
	http.HandleFunc("/getUser", addHeaders(ensureAuth(s.handleGetUser)))

	http.HandleFunc("/getConversations", addHeaders(ensureAuth(s.handleGetConversations)))
	http.HandleFunc("/getContacts", addHeaders(ensureAuth(s.handleGetContacts)))
	http.HandleFunc("/getConversation", addHeaders(ensureAuth(s.handleGetConversation)))
	http.HandleFunc("/getConversationsMessages", addHeaders(ensureAuth(s.handleGetConversationsMessages)))

	http.HandleFunc("/createConversation", addHeaders(ensureAuth(s.handleCreateConversation)))
	http.HandleFunc("/deleteConversation", addHeaders(ensureAuth(s.handleDeleteConversation)))
	http.HandleFunc("/joinConversation", addHeaders(ensureAuth(s.handleJoinPublicConversation)))
	http.HandleFunc("/leaveConversation", addHeaders(ensureAuth(s.handleLeavePublicConversation)))
	http.HandleFunc("/renameConversation", addHeaders(ensureAuth(s.handleRenamePublicConversation)))
}
