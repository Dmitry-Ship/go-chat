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

type eventHandlers struct {
	ctx        context.Context
	subscriber infra.EventsSubscriber
	commands   *app.Commands
	queries    readModel.QueriesRepository
}

func NewEventHandlers(ctx context.Context, subscriber infra.EventsSubscriber, commands *app.Commands, queries readModel.QueriesRepository) *eventHandlers {
	return &eventHandlers{
		ctx:        ctx,
		subscriber: subscriber,
		commands:   commands,
		queries:    queries,
	}
}

func (h *eventHandlers) logHandlerError(err error, e domain.DomainEvent) {
	h.logHandlerError(err, e)
}

func (h *eventHandlers) ListenForEvents() {
	for {
		select {
		case event := <-h.subscriber.Subscribe(domain.DomainEventChannel):
			switch e := event.Data.(type) {
			case *domain.PublicConversationRenamed:
				go h.handlePublicConversationRenamed(e)
			case *domain.PublicConversationLeft:
				go h.handlePublicConversationLeft(e)
			case *domain.PublicConversationJoined:
				go h.handlePublicConversationJoined(e)
			case *domain.MessageSent:
				go h.handleMessageSent(e)
			case *domain.PublicConversationCreated:
				go h.handlePublicConversationCreated(e)
			case *domain.PrivateConversationCreated:
				go h.handlePrivateConversationCreated(e)
			case *domain.PublicConversationDeleted:
				go h.handlePublicConversationDeleted(e)
			default:
				log.Printf("Unknown domain event type: %T", e)
			}

		case <-h.ctx.Done():
			return
		}
	}
}

func (h *eventHandlers) handlePublicConversationRenamed(e *domain.PublicConversationRenamed) {
	err := h.commands.MessagingService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

	if err != nil {
		h.logHandlerError(err, e)
	}

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
	err = h.commands.NotificationTopicService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

	if err != nil {
		h.logHandlerError(err, e)
	}
}

func (h *eventHandlers) handlePublicConversationLeft(e *domain.PublicConversationLeft) {
	err := h.commands.MessagingService.SendLeftConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err, e)
		return
	}

	err = h.commands.NotificationTopicService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		h.logHandlerError(err, e)
	}
}

func (h *eventHandlers) handlePublicConversationJoined(e *domain.PublicConversationJoined) {
	err := h.commands.MessagingService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		h.logHandlerError(err, e)
	}

	err = h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		h.logHandlerError(err, e)
	}
}

func (h *eventHandlers) handlePublicConversationDeleted(e *domain.PublicConversationDeleted) {
	notification := ws.OutgoingNotification{
		Type: "conversation_deleted",
		Payload: struct {
			ConversationId uuid.UUID `json:"conversation_id"`
		}{
			ConversationId: e.ConversationID,
		},
	}

	err := h.commands.NotificationTopicService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

	if err != nil {
		h.logHandlerError(err, e)
	}

	err = h.commands.NotificationTopicService.DeleteTopic("conversation:" + e.ConversationID.String())

	if err != nil {
		h.logHandlerError(err, e)
	}
}

func (h *eventHandlers) handleMessageSent(e *domain.MessageSent) {
	messageDTO, err := h.queries.GetNotificationMessage(e.MessageID, e.UserID)

	if err != nil {
		h.logHandlerError(err, e)
		return
	}

	notification := ws.OutgoingNotification{
		Type:    "message",
		Payload: messageDTO,
	}

	err = h.commands.NotificationTopicService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

	if err != nil {
		h.logHandlerError(err, e)
	}
}

func (h *eventHandlers) handlePublicConversationCreated(e *domain.PublicConversationCreated) {
	err := h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.OwnerID)

	if err != nil {
		h.logHandlerError(err, e)
	}
}

func (h *eventHandlers) handlePrivateConversationCreated(e *domain.PrivateConversationCreated) {
	err := h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.FromUserID)

	if err != nil {
		h.logHandlerError(err, e)
	}

	err = h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.ToUserID)

	if err != nil {
		h.logHandlerError(err, e)
	}
}
