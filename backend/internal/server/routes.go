package server

import (
	"net/http"
)

func (s *Server) initRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", corsHandler(s.config.ClientOrigin, securityHeaders(s.httpRateLimit(limitRequestBodySize(1<<20, s.handleSignUp)))))
	mux.HandleFunc("POST /login", corsHandler(s.config.ClientOrigin, securityHeaders(s.httpRateLimit(limitRequestBodySize(1<<20, s.handleLogin)))))
	mux.HandleFunc("POST /refreshToken", corsHandler(s.config.ClientOrigin, securityHeaders(s.httpRateLimit(limitRequestBodySize(1<<20, s.handleRefreshToken)))))
	mux.HandleFunc("POST /logout", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleLogout)))))

	mux.HandleFunc("GET /ws", securityHeaders(s.wsRateLimit(s.private(s.handleOpenWSConnection()))))

	mux.HandleFunc("POST /createConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleCreateGroupConversation)))))
	mux.HandleFunc("POST /startDirectConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleStartDirectConversation)))))
	mux.HandleFunc("POST /deleteConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleDeleteConversation)))))
	mux.HandleFunc("POST /joinConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleJoin)))))
	mux.HandleFunc("POST /inviteUserToConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleInvite)))))
	mux.HandleFunc("POST /kick", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleKick)))))
	mux.HandleFunc("POST /leaveConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleLeave)))))
	mux.HandleFunc("POST /renameConversation", corsHandler(s.config.ClientOrigin, securityHeaders(limitRequestBodySize(1<<20, s.private(s.handleRename)))))

	mux.HandleFunc("GET /getUser", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(s.handleGetUser))))
	mux.HandleFunc("GET /getConversations", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(withPagination(s.handleGetConversations)))))
	mux.HandleFunc("GET /getContacts", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(withPagination(s.handleGetContacts)))))
	mux.HandleFunc("GET /getPotentialInvitees", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(withPagination(s.handleGetPotentialInvitees)))))
	mux.HandleFunc("GET /getConversation", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(s.handleGetConversation))))
	mux.HandleFunc("GET /getConversationsMessages", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(withPagination(s.handleGetConversationsMessages)))))
	mux.HandleFunc("GET /getParticipants", corsHandler(s.config.ClientOrigin, securityHeaders(s.private(withPagination(s.handleGetParticipants)))))

	return mux
}
