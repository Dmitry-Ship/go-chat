package httpServer

import (
	ws "GitHub/go-chat/backend/pkg/websocket"
	"net/http"
)

func (s *HTTPServer) Init() {
	wsHandlers := ws.NewWSHandlers()
	wsHandlers.SetWSHandler("message", s.handleReceiveWSChatMessage)

	http.HandleFunc("/ws", s.private(s.handleOpenWSConnection(wsHandlers)))

	http.HandleFunc("/signup", s.withHeaders(s.handleSignUp))
	http.HandleFunc("/login", s.withHeaders(s.handleLogin))
	http.HandleFunc("/refreshToken", s.withHeaders((s.handleRefreshToken)))
	http.HandleFunc("/logout", s.private(s.handleLogout))

	http.HandleFunc("/getUser", s.private(s.handleGetUser))
	http.HandleFunc("/getConversations", s.private(s.handleGetConversations))
	http.HandleFunc("/getContacts", s.private(s.handleGetContacts))
	http.HandleFunc("/getConversation", s.private(s.handleGetConversation))
	http.HandleFunc("/getConversationsMessages", s.private(s.handleGetConversationsMessages))

	http.HandleFunc("/createConversation", s.private(s.handleCreateConversation))
	http.HandleFunc("/deleteConversation", s.private(s.handleDeleteConversation))
	http.HandleFunc("/joinConversation", s.private(s.handleJoinPublicConversation))
	http.HandleFunc("/leaveConversation", s.private(s.handleLeavePublicConversation))
	http.HandleFunc("/renameConversation", s.private(s.handleRenamePublicConversation))
}
