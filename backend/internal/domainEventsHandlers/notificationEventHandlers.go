package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"log"

	"github.com/google/uuid"
)

type notificationsEventHandlers struct {
	ctx                      context.Context
	subscriber               infra.EventsSubscriber
	notificationTopicService services.NotificationTopicService
	queries                  readModel.QueriesRepository
}

func NewNotificationsEventHandlers(ctx context.Context, subscriber infra.EventsSubscriber, notificationTopicService services.NotificationTopicService, queries readModel.QueriesRepository) *notificationsEventHandlers {
	return &notificationsEventHandlers{
		ctx:                      ctx,
		subscriber:               subscriber,
		notificationTopicService: notificationTopicService,
		queries:                  queries,
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
			case *domain.PublicConversationRenamed:
				go h.sendRenamedConversationNotification(e)
			case *domain.PublicConversationLeft:
				go h.unsubscribeFromConversation(e)
			case *domain.PublicConversationJoined:
				go h.subscribeToConversationNotifications(e.ConversationID, e.UserID)
			case *domain.PublicConversationInvited:
				go h.subscribeToConversationNotifications(e.ConversationID, e.UserID)
			case *domain.MessageSent:
				go h.sendMessageNotification(e)
			case *domain.PublicConversationCreated:
				go h.subscribeToConversationNotifications(e.ConversationID, e.OwnerID)
			case *domain.PrivateConversationCreated:
				go h.subscribeToConversationNotifications(e.ConversationID, e.FromUserID)
				go h.subscribeToConversationNotifications(e.ConversationID, e.ToUserID)
			case *domain.PublicConversationDeleted:
				go h.sendPublicConversationDeletedNotification(e)
				go h.deleteConversationTopic(e)
			}

		case <-h.ctx.Done():
			return
		}
	}
}

func (h *notificationsEventHandlers) unsubscribeFromConversation(e *domain.PublicConversationLeft) {
	err := h.notificationTopicService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *notificationsEventHandlers) deleteConversationTopic(e *domain.PublicConversationDeleted) {
	err := h.notificationTopicService.DeleteTopic("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *notificationsEventHandlers) subscribeToConversationNotifications(conversationId uuid.UUID, userID uuid.UUID) {
	err := h.notificationTopicService.SubscribeToTopic("conversation:"+conversationId.String(), userID)

	if err != nil {
		h.logHandlerError(err)
	}
}

func (h *notificationsEventHandlers) sendPublicConversationDeletedNotification(e *domain.PublicConversationDeleted) {
	ids, err := h.notificationTopicService.GetReceivers("conversation:" + e.ConversationID.String())

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

		err := h.notificationTopicService.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}

func (h *notificationsEventHandlers) sendRenamedConversationNotification(e *domain.PublicConversationRenamed) {
	ids, err := h.notificationTopicService.GetReceivers("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err)
		return
	}

	for _, id := range ids {
		notification := ws.OutgoingNotification{
			Type: "conversation_renamed",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
				NewName        string    `json:"new_name"`
			}{
				ConversationId: e.ConversationID,
				NewName:        e.NewName,
			},
		}

		err := h.notificationTopicService.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}

func (h *notificationsEventHandlers) sendMessageNotification(e *domain.MessageSent) {
	ids, err := h.notificationTopicService.GetReceivers("conversation:" + e.ConversationID.String())

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

		err = h.notificationTopicService.SendToUser(id, notification)

		if err != nil {
			h.logHandlerError(err)
		}
	}
}
