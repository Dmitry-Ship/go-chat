package server

import (
	"GitHub/go-chat/backend/internal/domain"
)

func (h *Server) sendRenamedConversationMessage(e *domain.GroupConversationRenamed) {
	err := h.conversationCommands.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *Server) sendGroupConversationLeftMessage(e *domain.GroupConversationLeft) {
	err := h.conversationCommands.SendLeftConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
		return
	}
}

func (h *Server) sendGroupConversationJoinedMessage(e *domain.GroupConversationJoined) {
	err := h.conversationCommands.SendJoinedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *Server) sendGroupConversationInvitedMessage(e *domain.GroupConversationInvited) {
	err := h.conversationCommands.SendInvitedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}
