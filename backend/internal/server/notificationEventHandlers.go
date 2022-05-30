package server

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (h *Server) sendGroupConversationDeletedNotification(e *domain.GroupConversationDeleted) {
	err := h.notificationCommands.SendToConversation(e.ConversationID, func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		notification := ws.OutgoingNotification{
			Type: "conversation_deleted",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: e.ConversationID,
			},
		}

		return &notification, nil
	})

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *Server) sendUpdatedConversationNotification(conversationID uuid.UUID) {
	err := h.notificationCommands.SendToConversation(conversationID, func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		conversation, err := h.queries.GetConversation(conversationID, userID)

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
		h.logHandlerError(err)
	}
}

func (h *Server) sendMessageNotification(e *domain.MessageSent) {
	err := h.notificationCommands.SendToConversation(e.ConversationID, func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
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
		h.logHandlerError(err)
	}
}
