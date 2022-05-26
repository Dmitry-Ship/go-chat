package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (h *eventHandlers) unsubscribeFromConversation(e *domain.GroupConversationLeft) {
	err := h.commands.NotificationService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *eventHandlers) deleteConversationTopic(e *domain.GroupConversationDeleted) {
	err := h.commands.NotificationService.DeleteTopic("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *eventHandlers) subscribeToConversationNotifications(conversationID uuid.UUID, userID uuid.UUID) {
	err := h.commands.NotificationService.SubscribeToTopic("conversation:"+conversationID.String(), userID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *eventHandlers) sendGroupConversationDeletedNotification(e *domain.GroupConversationDeleted) {
	ids, err := h.commands.NotificationService.GetReceivers("conversation:" + e.ConversationID.String())

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

		err := h.commands.NotificationService.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}

func (h *eventHandlers) sendUpdatedConversationNotification(conversationID uuid.UUID) {
	ids, err := h.commands.NotificationService.GetReceivers("conversation:" + conversationID.String())

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

		err = h.commands.NotificationService.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}

func (h *eventHandlers) sendMessageNotification(e *domain.MessageSent) {
	ids, err := h.commands.NotificationService.GetReceivers("conversation:" + e.ConversationID.String())

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

		err = h.commands.NotificationService.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}
