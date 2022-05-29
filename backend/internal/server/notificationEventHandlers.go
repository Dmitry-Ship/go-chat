package server

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (h *Server) unsubscribeFromConversation(e *domain.GroupConversationLeft) {
	err := h.notificationCommands.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *Server) deleteConversationTopic(e *domain.GroupConversationDeleted) {
	err := h.notificationCommands.DeleteTopic("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *Server) subscribeToConversationNotifications(conversationID uuid.UUID, userID uuid.UUID) {
	err := h.notificationCommands.SubscribeToTopic("conversation:"+conversationID.String(), userID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *Server) sendGroupConversationDeletedNotification(e *domain.GroupConversationDeleted) {
	ids, err := h.notificationCommands.GetReceivers("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
		return
	}

	for _, id := range ids {
		notification := ws.OutgoingNotification{
			Type: "conversation_deleted",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: e.ConversationID,
			},
		}

		err := h.notificationCommands.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}

func (h *Server) sendUpdatedConversationNotification(conversationID uuid.UUID) {
	ids, err := h.notificationCommands.GetReceivers("conversation:" + conversationID.String())

	if err != nil {
		h.logHandlerError(err)
		return
	}

	for _, id := range ids {
		conversation, err := h.queries.GetConversation(conversationID, id)

		if err != nil {
			h.logHandlerError(err)
			return
		}

		notification := ws.OutgoingNotification{
			Type:    "conversation_updated",
			Payload: conversation,
		}

		err = h.notificationCommands.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}

func (h *Server) sendMessageNotification(e *domain.MessageSent) {
	ids, err := h.notificationCommands.GetReceivers("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
		return
	}

	for _, id := range ids {
		messageDTO, err := h.queries.GetNotificationMessage(e.MessageID, id)

		if err != nil {
			h.logHandlerError(err)
			return
		}

		notification := ws.OutgoingNotification{
			Type:    "message",
			Payload: messageDTO,
		}

		err = h.notificationCommands.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}
