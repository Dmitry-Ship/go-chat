package server

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

func (h *Server) sendMessage(event domain.DomainEvent) error {
	switch e := event.(type) {
	case *domain.GroupConversationRenamed:
		return h.conversationCommands.SendRenamedConversationMessage(e.GetConversationID(), e.UserID, e.NewName)
	case *domain.GroupConversationLeft:
		return h.conversationCommands.SendLeftConversationMessage(e.GetConversationID(), e.UserID)
	case *domain.GroupConversationJoined:
		return h.conversationCommands.SendJoinedConversationMessage(e.GetConversationID(), e.UserID)
	case *domain.GroupConversationInvited:
		return h.conversationCommands.SendInvitedConversationMessage(e.GetConversationID(), e.UserID)
	}

	return nil
}

func (h *Server) sendWSNotification(event domain.DomainEvent) error {
	var targetIDs []uuid.UUID
	var err error
	var buildMessage func(userID uuid.UUID) (*ws.OutgoingNotification, error)

	switch e := event.(type) {
	case *domain.GroupConversationRenamed:
		targetIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())
		buildMessage = h.notificationBuilder.GetConversationUpdatedBuilder(e.GetConversationID())
	case *domain.GroupConversationLeft:
		targetIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())
		buildMessage = h.notificationBuilder.GetConversationUpdatedBuilder(e.GetConversationID())
	case *domain.GroupConversationJoined:
		targetIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())
		buildMessage = h.notificationBuilder.GetConversationUpdatedBuilder(e.GetConversationID())
	case *domain.GroupConversationInvited:
		targetIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())
		buildMessage = h.notificationBuilder.GetConversationUpdatedBuilder(e.GetConversationID())
	case *domain.GroupConversationDeleted:
		targetIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())
		buildMessage = h.notificationBuilder.GetConversationDeletedBuilder(e.GetConversationID())
	case *domain.MessageSent:
		targetIDs, err = h.notificationResolver.GetConversationRecipients(e.GetConversationID())
		buildMessage = h.notificationBuilder.GetMessageSentBuilder(e.MessageID)
	}

	if err != nil {
		return err
	}

	return h.notificationCommands.Broadcast(targetIDs, buildMessage)
}
