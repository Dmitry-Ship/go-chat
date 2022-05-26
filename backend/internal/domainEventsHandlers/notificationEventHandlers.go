package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"log"

	"github.com/google/uuid"
)

type notificationsEventHandlers struct {
	ctx        context.Context
	subscriber infra.EventsSubscriber
	commands   *app.Commands
	queries    readModel.QueriesRepository
}

func NewNotificationsEventHandlers(ctx context.Context, subscriber infra.EventsSubscriber, commands *app.Commands, queries readModel.QueriesRepository) *notificationsEventHandlers {
	return &notificationsEventHandlers{
		ctx:        ctx,
		subscriber: subscriber,
		commands:   commands,
		queries:    queries,
	}
}

func (h *notificationsEventHandlers) logHandlerError(err error) {
	log.Panicln("Error occurred event", err)
}

func (h *notificationsEventHandlers) ListenForEvents() {
	for {
		select {
		case event := <-h.subscriber.Subscribe(domain.DomainEventChannel):
			switch e := event.Data.(type) {
			case *domain.GroupConversationRenamed:
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationLeft:
				go h.unsubscribeFromConversation(e)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationJoined:
				go h.subscribeToConversationNotifications(e.ConversationID, e.UserID)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.GroupConversationInvited:
				go h.subscribeToConversationNotifications(e.ConversationID, e.UserID)
				go h.sendUpdatedConversationNotification(e.ConversationID)
			case *domain.MessageSent:
				go h.sendMessageNotification(e)
			case *domain.GroupConversationCreated:
				go h.subscribeToConversationNotifications(e.ConversationID, e.OwnerID)
			case *domain.DirectConversationCreated:
				go h.subscribeToConversationNotifications(e.ConversationID, e.FromUserID)
				go h.subscribeToConversationNotifications(e.ConversationID, e.ToUserID)
			case *domain.GroupConversationDeleted:
				go h.sendGroupConversationDeletedNotification(e)
				go h.deleteConversationTopic(e)
			}

		case <-h.ctx.Done():
			return
		}
	}
}

func (h *notificationsEventHandlers) unsubscribeFromConversation(e *domain.GroupConversationLeft) {
	err := h.commands.NotificationService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *notificationsEventHandlers) deleteConversationTopic(e *domain.GroupConversationDeleted) {
	err := h.commands.NotificationService.DeleteTopic("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *notificationsEventHandlers) subscribeToConversationNotifications(conversationID uuid.UUID, userID uuid.UUID) {
	err := h.commands.NotificationService.SubscribeToTopic("conversation:"+conversationID.String(), userID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *notificationsEventHandlers) sendGroupConversationDeletedNotification(e *domain.GroupConversationDeleted) {
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

func (h *notificationsEventHandlers) sendUpdatedConversationNotification(conversationID uuid.UUID) {
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

func (h *notificationsEventHandlers) sendMessageNotification(e *domain.MessageSent) {
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
