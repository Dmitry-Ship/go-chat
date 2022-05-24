package httpHandlers

import (
	"net/http"
)

func (s *HTTPHandlers) InitRoutes() {
	http.HandleFunc("/ws", s.direct(s.commandController.handleOpenWSConnection()))

	http.HandleFunc("/signup", s.Post(s.withHeaders(s.commandController.handleSignUp)))
	http.HandleFunc("/login", s.Post(s.withHeaders(s.commandController.handleLogin)))
	http.HandleFunc("/refreshToken", s.Post(s.withHeaders((s.commandController.handleRefreshToken))))
	http.HandleFunc("/logout", s.direct(s.Post(s.commandController.handleLogout)))

	http.HandleFunc("/getUser", s.direct(s.Get(s.queryController.handleGetUser)))
	http.HandleFunc("/getConversations", s.direct(s.GetPaginated(s.queryController.handleGetConversations)))
	http.HandleFunc("/getContacts", s.direct(s.GetPaginated(s.queryController.handleGetContacts)))
	http.HandleFunc("/getPotentialInvitees", s.direct(s.GetPaginated(s.queryController.handleGetPotentialInvitees)))
	http.HandleFunc("/getConversation", s.direct(s.Get(s.queryController.handleGetConversation)))
	http.HandleFunc("/getConversationsMessages", (s.direct(s.GetPaginated(s.queryController.handleGetConversationsMessages))))
	http.HandleFunc("/getParticipants", s.direct(s.GetPaginated(s.queryController.handleGetParticipants)))

	http.HandleFunc("/createConversation", s.direct(s.Post(s.commandController.handleCreateGroupConversation)))
	http.HandleFunc("/createDirectConversationIfNotExists", s.direct(s.Post(s.commandController.handleCreateDirectConversationIfNotExists)))
	http.HandleFunc("/deleteConversation", s.direct(s.Post(s.commandController.handleDeleteConversation)))
	http.HandleFunc("/joinConversation", s.direct(s.Post(s.commandController.handleJoinGroupConversation)))
	http.HandleFunc("/inviteUserToConversation", s.direct(s.Post(s.commandController.handleInviteToGroupConversation)))
	http.HandleFunc("/leaveConversation", s.direct(s.Post(s.commandController.handleLeaveGroupConversation)))
	http.HandleFunc("/renameConversation", s.direct(s.Post(s.commandController.handleRenameGroupConversation)))
}
