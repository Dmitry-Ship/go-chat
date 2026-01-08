package server

import (
	"net/http"
)

func (s *Server) initRoutes() {
	http.HandleFunc("/signup", securityHeaders(s.httpRateLimit(limitRequestBodySize(1<<20, s.post(s.handleSignUp)))))
	http.HandleFunc("/login", securityHeaders(s.httpRateLimit(limitRequestBodySize(1<<20, s.post(s.handleLogin)))))
	http.HandleFunc("/refreshToken", securityHeaders(s.httpRateLimit(limitRequestBodySize(1<<20, s.post(s.handleRefreshToken)))))
	http.HandleFunc("/logout", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleLogout)))))

	http.HandleFunc("/ws", securityHeaders(s.wsRateLimit(s.private(s.handleOpenWSConnection()))))

	http.HandleFunc("/createConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleCreateGroupConversation)))))
	http.HandleFunc("/startDirectConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleStartDirectConversation)))))
	http.HandleFunc("/deleteConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleDeleteConversation)))))
	http.HandleFunc("/joinConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleJoin)))))
	http.HandleFunc("/inviteUserToConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleInvite)))))
	http.HandleFunc("/kick", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleKick)))))
	http.HandleFunc("/leaveConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleLeave)))))
	http.HandleFunc("/renameConversation", securityHeaders(limitRequestBodySize(1<<20, s.post(s.private(s.handleRename)))))

	http.HandleFunc("/getUser", securityHeaders(s.get(s.private(s.handleGetUser))))
	http.HandleFunc("/getConversations", securityHeaders(s.getPaginated(s.private(s.handleGetConversations))))
	http.HandleFunc("/getContacts", securityHeaders(s.getPaginated(s.private(s.handleGetContacts))))
	http.HandleFunc("/getPotentialInvitees", securityHeaders(s.getPaginated(s.private(s.handleGetPotentialInvitees))))
	http.HandleFunc("/getConversation", securityHeaders(s.get(s.private(s.handleGetConversation))))
	http.HandleFunc("/getConversationsMessages", securityHeaders(s.getPaginated(s.private(s.handleGetConversationsMessages))))
	http.HandleFunc("/getParticipants", securityHeaders(s.getPaginated(s.private(s.handleGetParticipants))))
}
