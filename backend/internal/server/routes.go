package server

import (
	"net/http"
)

func (s *Server) initRoutes() {
	http.HandleFunc("/signup", s.post(s.handleSignUp))
	http.HandleFunc("/login", s.post(s.handleLogin))
	http.HandleFunc("/refreshToken", s.post(s.handleRefreshToken))
	http.HandleFunc("/logout", s.post(s.private(s.handleLogout)))

	http.HandleFunc("/ws", s.wsRateLimit(s.private(s.handleOpenWSConnection())))

	http.HandleFunc("/createConversation", s.post(s.private(s.handleCreateGroupConversation)))
	http.HandleFunc("/startDirectConversation", s.post(s.private(s.handleStartDirectConversation)))
	http.HandleFunc("/deleteConversation", s.post(s.private(s.handleDeleteConversation)))
	http.HandleFunc("/joinConversation", s.post(s.private(s.handleJoin)))
	http.HandleFunc("/inviteUserToConversation", s.post(s.private(s.handleInvite)))
	http.HandleFunc("/kick", s.post(s.private(s.handleKick)))
	http.HandleFunc("/leaveConversation", s.post(s.private(s.handleLeave)))
	http.HandleFunc("/renameConversation", s.post(s.private(s.handleRename)))

	http.HandleFunc("/getUser", s.get(s.private(s.handleGetUser)))
	http.HandleFunc("/getConversations", s.getPaginated(s.private(s.handleGetConversations)))
	http.HandleFunc("/getContacts", s.getPaginated(s.private(s.handleGetContacts)))
	http.HandleFunc("/getPotentialInvitees", s.getPaginated(s.private(s.handleGetPotentialInvitees)))
	http.HandleFunc("/getConversation", s.get(s.private(s.handleGetConversation)))
	http.HandleFunc("/getConversationsMessages", s.getPaginated(s.private(s.handleGetConversationsMessages)))
	http.HandleFunc("/getParticipants", s.getPaginated(s.private(s.handleGetParticipants)))
}
