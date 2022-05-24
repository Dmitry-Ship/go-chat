package httpHandlers

import (
	"net/http"
)

func (s *HTTPHandlers) InitRoutes() {
	http.HandleFunc("/ws", s.private(s.commandController.handleOpenWSConnection()))

	http.HandleFunc("/signup", s.Post(s.withHeaders(s.commandController.handleSignUp)))
	http.HandleFunc("/login", s.Post(s.withHeaders(s.commandController.handleLogin)))
	http.HandleFunc("/refreshToken", s.Post(s.withHeaders((s.commandController.handleRefreshToken))))
	http.HandleFunc("/logout", s.private(s.Post(s.commandController.handleLogout)))

	http.HandleFunc("/getUser", s.private(s.Get(s.queryController.handleGetUser)))
	http.HandleFunc("/getConversations", s.private(s.GetPaginated(s.queryController.handleGetConversations)))
	http.HandleFunc("/getContacts", s.private(s.GetPaginated(s.queryController.handleGetContacts)))
	http.HandleFunc("/getPotentialInvitees", s.private(s.GetPaginated(s.queryController.handleGetPotentialInvitees)))
	http.HandleFunc("/getConversation", s.private(s.Get(s.queryController.handleGetConversation)))
	http.HandleFunc("/getConversationsMessages", (s.private(s.GetPaginated(s.queryController.handleGetConversationsMessages))))
	http.HandleFunc("/getParticipants", s.private(s.GetPaginated(s.queryController.handleGetParticipants)))

	http.HandleFunc("/createConversation", s.private(s.Post(s.commandController.handleCreatePublicConversation)))
	http.HandleFunc("/createPrivateConversationIfNotExists", s.private(s.Post(s.commandController.handleCreatePrivateConversationIfNotExists)))
	http.HandleFunc("/deleteConversation", s.private(s.Post(s.commandController.handleDeleteConversation)))
	http.HandleFunc("/joinConversation", s.private(s.Post(s.commandController.handleJoinPublicConversation)))
	http.HandleFunc("/inviteUserToConversation", s.private(s.Post(s.commandController.handleInviteToPublicConversation)))
	http.HandleFunc("/leaveConversation", s.private(s.Post(s.commandController.handleLeavePublicConversation)))
	http.HandleFunc("/renameConversation", s.private(s.Post(s.commandController.handleRenamePublicConversation)))
}
