package httpHandlers

import (
	"net/http"
)

func (s *HTTPHandlers) InitRoutes() {
	http.HandleFunc("/ws", s.private(s.commandController.handleOpenWSConnection()))

	http.HandleFunc("/signup", s.withHeaders(s.commandController.handleSignUp))
	http.HandleFunc("/login", s.withHeaders(s.commandController.handleLogin))
	http.HandleFunc("/refreshToken", s.withHeaders((s.commandController.handleRefreshToken)))
	http.HandleFunc("/logout", s.private(s.commandController.handleLogout))

	http.HandleFunc("/getUser", s.private(s.queryController.handleGetUser))
	http.HandleFunc("/getConversations", s.private(s.paginate(s.queryController.handleGetConversations)))
	http.HandleFunc("/getContacts", s.private(s.paginate(s.queryController.handleGetContacts)))
	http.HandleFunc("/getPotentialInvitees", s.private(s.paginate(s.queryController.handleGetPotentialInvitees)))
	http.HandleFunc("/getConversation", s.private(s.queryController.handleGetConversation))
	http.HandleFunc("/getConversationsMessages", s.private(s.paginate(s.queryController.handleGetConversationsMessages)))
	http.HandleFunc("/getParticipants", s.private(s.paginate(s.queryController.handleGetParticipants)))

	http.HandleFunc("/createConversation", s.private(s.commandController.handleCreatePublicConversation))
	http.HandleFunc("/createPrivateConversationIfNotExists", s.private(s.commandController.handleCreatePrivateConversationIfNotExists))
	http.HandleFunc("/deleteConversation", s.private(s.commandController.handleDeleteConversation))
	http.HandleFunc("/joinConversation", s.private(s.commandController.handleJoinPublicConversation))
	http.HandleFunc("/inviteUserToConversation", s.private(s.commandController.handleInviteToPublicConversation))
	http.HandleFunc("/leaveConversation", s.private(s.commandController.handleLeavePublicConversation))
	http.HandleFunc("/renameConversation", s.private(s.commandController.handleRenamePublicConversation))
}
