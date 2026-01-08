package server

import (
	"net/http"
)

func (s *Server) initRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", s.corsHandler(s.securityHeaders(s.httpRateLimit(s.limitRequestBodySize(MaxRequestBodySize, s.handleSignUp)))))
	mux.HandleFunc("POST /login", s.corsHandler(s.securityHeaders(s.httpRateLimit(s.limitRequestBodySize(MaxRequestBodySize, s.handleLogin)))))
	mux.HandleFunc("POST /refreshToken", s.corsHandler(s.securityHeaders(s.httpRateLimit(s.limitRequestBodySize(MaxRequestBodySize, s.handleRefreshToken)))))
	mux.HandleFunc("POST /logout", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleLogout)))))

	mux.HandleFunc("GET /ws", s.securityHeaders(s.wsRateLimit(s.private(s.handleOpenWSConnection()))))

	mux.HandleFunc("POST /createConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleCreateGroupConversation)))))
	mux.HandleFunc("POST /startDirectConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleStartDirectConversation)))))
	mux.HandleFunc("POST /deleteConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleDeleteConversation)))))
	mux.HandleFunc("POST /joinConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleJoin)))))
	mux.HandleFunc("POST /inviteUserToConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleInvite)))))
	mux.HandleFunc("POST /kick", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleKick)))))
	mux.HandleFunc("POST /leaveConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleLeave)))))
	mux.HandleFunc("POST /renameConversation", s.corsHandler(s.securityHeaders(s.limitRequestBodySize(MaxRequestBodySize, s.private(s.handleRename)))))

	mux.HandleFunc("GET /getUser", s.corsHandler(s.securityHeaders(s.private(s.handleGetUser))))
	mux.HandleFunc("GET /getConversations", s.corsHandler(s.securityHeaders(s.private(withPagination(s.handleGetConversations)))))
	mux.HandleFunc("GET /getContacts", s.corsHandler(s.securityHeaders(s.private(withPagination(s.handleGetContacts)))))
	mux.HandleFunc("GET /getPotentialInvitees", s.corsHandler(s.securityHeaders(s.private(withPagination(s.handleGetPotentialInvitees)))))
	mux.HandleFunc("GET /getConversation", s.corsHandler(s.securityHeaders(s.private(s.handleGetConversation))))
	mux.HandleFunc("GET /getConversationsMessages", s.corsHandler(s.securityHeaders(s.private(withPagination(s.handleGetConversationsMessages)))))
	mux.HandleFunc("GET /getParticipants", s.corsHandler(s.securityHeaders(s.private(withPagination(s.handleGetParticipants)))))

	return mux
}
