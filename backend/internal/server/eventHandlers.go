package server

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (h *Server) sendRenamedConversationMessage(e *domain.GroupConversationRenamed) error {
	return h.conversationCommands.SendRenamedConversationMessage(e.GetConversationID(), e.UserID, e.NewName)
}

func (h *Server) sendGroupConversationLeftMessage(e *domain.GroupConversationLeft) error {
	return h.conversationCommands.SendLeftConversationMessage(e.GetConversationID(), e.UserID)
}

func (h *Server) sendGroupConversationJoinedMessage(e *domain.GroupConversationJoined) error {
	return h.conversationCommands.SendJoinedConversationMessage(e.GetConversationID(), e.UserID)
}

func (h *Server) sendGroupConversationInvitedMessage(e *domain.GroupConversationInvited) error {
	return h.conversationCommands.SendInvitedConversationMessage(e.GetConversationID(), e.UserID)
}

func (h *Server) sendGroupConversationDeletedNotification(e *domain.GroupConversationDeleted) error {
	return h.notificationCommands.SendToConversation(e.GetConversationID(), func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		return &ws.OutgoingNotification{
			Type: "conversation_deleted",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: e.GetConversationID(),
			},
		}, nil
	})
}

func (h *Server) sendUpdatedConversationNotification(e domain.ConversationEvent) error {
	return h.notificationCommands.SendToConversation(e.GetConversationID(), func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		conversation, err := h.queries.GetConversation(e.GetConversationID(), userID)

		if err != nil {
			return nil, err
		}

		return &ws.OutgoingNotification{
			Type:    "conversation_updated",
			Payload: conversation,
		}, nil
	})
}

func (h *Server) sendMessageNotification(e *domain.MessageSent) error {
	return h.notificationCommands.SendToConversation(e.GetConversationID(), func(userID uuid.UUID) (*ws.OutgoingNotification, error) {
		messageDTO, err := h.queries.GetNotificationMessage(e.MessageID, userID)

		if err != nil {
			return nil, err
		}

		return &ws.OutgoingNotification{
			Type:    "message",
			Payload: messageDTO,
		}, nil
	})
}
