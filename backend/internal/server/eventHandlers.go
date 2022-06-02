package server

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (h *Server) sendRenamedConversationMessage(e *domain.GroupConversationRenamed) {
	err := h.conversationCommands.SendRenamedConversationMessage(e.GetConversationID(), e.UserID, e.NewName)

	if err != nil {
		h.logHandlerError(e, err)
	}
}

func (h *Server) sendGroupConversationLeftMessage(e *domain.GroupConversationLeft) {
	err := h.conversationCommands.SendLeftConversationMessage(e.GetConversationID(), e.UserID)

	if err != nil {
		h.logHandlerError(e, err)
		return
	}
}

func (h *Server) sendGroupConversationJoinedMessage(e *domain.GroupConversationJoined) {
	err := h.conversationCommands.SendJoinedConversationMessage(e.GetConversationID(), e.UserID)

	if err != nil {
		h.logHandlerError(e, err)
	}
}

func (h *Server) sendGroupConversationInvitedMessage(e *domain.GroupConversationInvited) {
	err := h.conversationCommands.SendInvitedConversationMessage(e.GetConversationID(), e.UserID)

	if err != nil {
		h.logHandlerError(e, err)
	}
}

func (h *Server) sendGroupConversationDeletedNotification(e *domain.GroupConversationDeleted) {
	err := h.notificationCommands.SendToConversation(e.GetConversationID(), func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		notification := ws.OutgoingNotification{
			Type: "conversation_deleted",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: e.GetConversationID(),
			},
		}

		return &notification, nil
	})

	if err != nil {
		h.logHandlerError(e, err)
	}
}

func (h *Server) sendUpdatedConversationNotification(e domain.ConversationEvent) {
	err := h.notificationCommands.SendToConversation(e.GetConversationID(), func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		conversation, err := h.queries.GetConversation(e.GetConversationID(), userID)

		if err != nil {
			return nil, err
		}

		notification := ws.OutgoingNotification{
			Type:    "conversation_updated",
			Payload: conversation,
		}

		return &notification, nil
	})

	if err != nil {
		h.logHandlerError(e, err)
	}
}

func (h *Server) sendMessageNotification(e *domain.MessageSent) {
	err := h.notificationCommands.SendToConversation(e.GetConversationID(), func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		messageDTO, err := h.queries.GetNotificationMessage(e.MessageID, userID)

		if err != nil {
			return nil, err
		}

		notification := ws.OutgoingNotification{
			Type:    "message",
			Payload: messageDTO,
		}

		return &notification, nil
	})

	if err != nil {
		h.logHandlerError(e, err)
	}
}
