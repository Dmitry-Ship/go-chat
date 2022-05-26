package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/domain"
)

func (h *eventHandlers) sendRenamedConversationMessage(e *domain.GroupConversationRenamed) {
	err := h.commands.ConversationService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *eventHandlers) sendGroupConversationLeftMessage(e *domain.GroupConversationLeft) {
	err := h.commands.ConversationService.SendLeftConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
		return
	}
}

func (h *eventHandlers) sendGroupConversationJoinedMessage(e *domain.GroupConversationJoined) {
	err := h.commands.ConversationService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *eventHandlers) sendGroupConversationInvitedMessage(e *domain.GroupConversationInvited) {
	err := h.commands.ConversationService.SendInvitedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}
