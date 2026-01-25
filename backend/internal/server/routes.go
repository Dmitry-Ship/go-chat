package server

import (
	"net/http"
)

func (s *Server) initRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("OPTIONS /api/{rest...}", s.securityHeaders(func(w http.ResponseWriter, r *http.Request) {}))
	mux.HandleFunc("POST /api/signup", s.securityHeaders(s.httpRateLimit(s.limitRequestBodySize(MaxRequestBodySize, s.handleSignUp))))
	mux.HandleFunc("POST /api/login", s.securityHeaders(s.httpRateLimit(s.limitRequestBodySize(MaxRequestBodySize, s.handleLogin))))
	mux.HandleFunc("POST /api/refreshToken", s.securityHeaders(s.httpRateLimit(s.handleRefreshToken)))
	mux.HandleFunc("POST /api/logout", s.securityHeaders(s.private(s.handleLogout)))

	mux.HandleFunc("GET /ws", s.securityHeaders(s.wsRateLimit(s.private(s.handleOpenWSConnection()))))

	mux.HandleFunc("POST /api/createConversation", s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleCreateGroupConversation))))
	mux.HandleFunc("POST /api/sendMessage", s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleSendMessage))))
	mux.HandleFunc("POST /api/startDirectConversation", s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleStartDirectConversation))))
	mux.HandleFunc("POST /api/deleteConversation", s.securityHeaders(s.private(s.handleDeleteConversation)))
	mux.HandleFunc("POST /api/joinConversation", s.securityHeaders(s.private(s.handleJoin)))
	mux.HandleFunc("POST /api/inviteUserToConversation", s.securityHeaders(s.private(s.handleInvite)))
	mux.HandleFunc("POST /api/kick", s.securityHeaders(s.private(s.handleKick)))
	mux.HandleFunc("POST /api/leaveConversation", s.securityHeaders(s.private(s.handleLeave)))
	mux.HandleFunc("POST /api/renameConversation", s.securityHeaders(s.private(s.handleRename)))

	mux.HandleFunc("GET /api/getUser", s.securityHeaders(s.private(s.handleGetUser)))
	mux.HandleFunc("GET /api/getConversations", s.securityHeaders(s.private(withPagination(s.handleGetConversations))))
	mux.HandleFunc("GET /api/getContacts", s.securityHeaders(s.private(withPagination(s.handleGetContacts))))
	mux.HandleFunc("GET /api/getPotentialInvitees", s.securityHeaders(s.private(withPagination(s.handleGetPotentialInvitees))))
	mux.HandleFunc("GET /api/getConversation", s.securityHeaders(s.private(s.handleGetConversation)))
	mux.HandleFunc("GET /api/getConversationsMessages", s.securityHeaders(s.private(s.handleGetConversationsMessages)))
	mux.HandleFunc("GET /api/getConversationUsers", s.securityHeaders(s.private(s.handleGetConversationUsers)))
	mux.HandleFunc("GET /api/getParticipants", s.securityHeaders(s.private(withPagination(s.handleGetParticipants))))

	return mux
}
