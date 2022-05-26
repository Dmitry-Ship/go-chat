package httpHandlers

import (
	"net/http"
)

func (s *HTTPHandlers) InitRoutes() {
	http.HandleFunc("/signup", s.post(s.commandHandlers.handleSignUp))
	http.HandleFunc("/login", s.post(s.commandHandlers.handleLogin))
	http.HandleFunc("/refreshToken", s.post(s.commandHandlers.handleRefreshToken))
	http.HandleFunc("/logout", s.post(s.private(s.commandHandlers.handleLogout)))

	http.HandleFunc("/ws", s.private(s.commandHandlers.handleOpenWSConnection()))

	http.HandleFunc("/createConversation", s.post(s.private(s.commandHandlers.handleCreateGroupConversation)))
	http.HandleFunc("/startDirectConversation", s.post(s.private(s.commandHandlers.handleStartDirectConversation)))
	http.HandleFunc("/deleteConversation", s.post(s.private(s.commandHandlers.handleDeleteConversation)))
	http.HandleFunc("/joinConversation", s.post(s.private(s.commandHandlers.handleJoinGroupConversation)))
	http.HandleFunc("/inviteUserToConversation", s.post(s.private(s.commandHandlers.handleInviteToGroupConversation)))
	http.HandleFunc("/leaveConversation", s.post(s.private(s.commandHandlers.handleLeaveGroupConversation)))
	http.HandleFunc("/renameConversation", s.post(s.private(s.commandHandlers.handleRenameGroupConversation)))

	http.HandleFunc("/getUser", s.get(s.private(s.queryHandlers.handleGetUser)))
	http.HandleFunc("/getConversations", s.getPaginated(s.private(s.queryHandlers.handleGetConversations)))
	http.HandleFunc("/getContacts", s.getPaginated(s.private(s.queryHandlers.handleGetContacts)))
	http.HandleFunc("/getPotentialInvitees", s.getPaginated(s.private(s.queryHandlers.handleGetPotentialInvitees)))
	http.HandleFunc("/getConversation", s.get(s.private(s.queryHandlers.handleGetConversation)))
	http.HandleFunc("/getConversationsMessages", s.getPaginated(s.private(s.queryHandlers.handleGetConversationsMessages)))
	http.HandleFunc("/getParticipants", s.getPaginated(s.private(s.queryHandlers.handleGetParticipants)))
}
