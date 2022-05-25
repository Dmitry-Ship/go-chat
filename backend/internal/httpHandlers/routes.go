package httpHandlers

import (
	"net/http"
)

func (s *HTTPHandlers) InitRoutes() {
	http.HandleFunc("/ws", s.private(s.commandController.handleOpenWSConnection()))

	http.HandleFunc("/signup", s.post(s.commandController.handleSignUp))
	http.HandleFunc("/login", s.post(s.commandController.handleLogin))
	http.HandleFunc("/refreshToken", s.post(s.commandController.handleRefreshToken))
	http.HandleFunc("/logout", s.post(s.private(s.commandController.handleLogout)))

	http.HandleFunc("/getUser", s.get(s.private(s.queryController.handleGetUser)))
	http.HandleFunc("/getConversations", s.getPaginated(s.private(s.queryController.handleGetConversations)))
	http.HandleFunc("/getContacts", s.getPaginated(s.private(s.queryController.handleGetContacts)))
	http.HandleFunc("/getPotentialInvitees", s.getPaginated(s.private(s.queryController.handleGetPotentialInvitees)))
	http.HandleFunc("/getConversation", s.get(s.private(s.queryController.handleGetConversation)))
	http.HandleFunc("/getConversationsMessages", s.getPaginated(s.private(s.queryController.handleGetConversationsMessages)))
	http.HandleFunc("/getParticipants", s.getPaginated(s.private(s.queryController.handleGetParticipants)))

	http.HandleFunc("/createConversation", s.post(s.private(s.commandController.handleCreateGroupConversation)))
	http.HandleFunc("/createDirectConversationIfNotExists", s.post(s.private(s.commandController.handleCreateDirectConversationIfNotExists)))
	http.HandleFunc("/deleteConversation", s.post(s.private(s.commandController.handleDeleteConversation)))
	http.HandleFunc("/joinConversation", s.post(s.private(s.commandController.handleJoinGroupConversation)))
	http.HandleFunc("/inviteUserToConversation", s.post(s.private(s.commandController.handleInviteToGroupConversation)))
	http.HandleFunc("/leaveConversation", s.post(s.private(s.commandController.handleLeaveGroupConversation)))
	http.HandleFunc("/renameConversation", s.post(s.private(s.commandController.handleRenameGroupConversation)))
}
