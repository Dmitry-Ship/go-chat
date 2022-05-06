package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"
	"log"

	"github.com/google/uuid"
)

type eventHandlers struct {
	ctx        context.Context
	subscriber domain.EventsSubscriber
	commands   *app.Commands
	queries    readModel.QueriesRepository
}

func NewEventHandlers(ctx context.Context, subscriber domain.EventsSubscriber, commands *app.Commands, queries readModel.QueriesRepository) *eventHandlers {
	return &eventHandlers{
		ctx:        ctx,
		subscriber: subscriber,
		commands:   commands,
		queries:    queries,
	}
}

func (h *eventHandlers) ListenForEvents() {
	for {
		select {
		case event := <-h.subscriber.Subscribe(domain.PublicConversationRenamedName):
			go h.handlePublicConversationRenamed(event)
		case event := <-h.subscriber.Subscribe(domain.PublicConversationLeftName):
			go h.handlePublicConversationLeft(event)
		case event := <-h.subscriber.Subscribe(domain.PublicConversationJoinedName):
			go h.handlePublicConversationJoined(event)
		case event := <-h.subscriber.Subscribe(domain.MessageSentName):
			go h.handleMessageSent(event)
		case event := <-h.subscriber.Subscribe(domain.PublicConversationCreatedName):
			go h.handlePublicConversationCreated(event)
		case event := <-h.subscriber.Subscribe(domain.PrivateConversationCreatedName):
			go h.handlePrivateConversationCreated(event)
		case event := <-h.subscriber.Subscribe(domain.PublicConversationDeletedName):
			go h.handlePublicConversationDeleted(event)
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *eventHandlers) handlePublicConversationRenamed(event domain.DomainEvent) {
	e, ok := event.(*domain.PublicConversationRenamed)
	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

	err := h.commands.MessagingService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
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
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}

func (h *eventHandlers) handlePublicConversationLeft(event domain.DomainEvent) {
	e, ok := event.(*domain.PublicConversationLeft)

	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

	err := h.commands.MessagingService.SendLeftConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}

	err = h.commands.NotificationTopicService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}

func (h *eventHandlers) handlePublicConversationJoined(event domain.DomainEvent) {
	e, ok := event.(*domain.PublicConversationJoined)

	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

	err := h.commands.MessagingService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}

	err = h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.UserID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}

func (h *eventHandlers) handlePublicConversationDeleted(event domain.DomainEvent) {
	e, ok := event.(*domain.PublicConversationDeleted)

	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

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
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}

	err = h.commands.NotificationTopicService.DeleteTopic("conversation:" + e.ConversationID.String())

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}

func (h *eventHandlers) handleMessageSent(event domain.DomainEvent) {
	e, ok := event.(*domain.MessageSent)

	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

	messageDTO, err := h.queries.GetNotificationMessage(e.MessageID, e.UserID)

	if err != nil {
		log.Println(err)
		return
	}

	notification := ws.OutgoingNotification{
		Type:    "message",
		Payload: messageDTO,
	}

	err = h.commands.NotificationTopicService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}

func (h *eventHandlers) handlePublicConversationCreated(event domain.DomainEvent) {
	e, ok := event.(*domain.PublicConversationCreated)

	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

	err := h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.OwnerID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}

func (h *eventHandlers) handlePrivateConversationCreated(event domain.DomainEvent) {
	e, ok := event.(*domain.PrivateConversationCreated)

	if !ok {
		log.Printf("Wrong event type: %T", e)
		return
	}

	err := h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.FromUserID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}

	err = h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.ToUserID)

	if err != nil {
		log.Printf("Error handling %s event: %v", e.GetName(), err)
	}
}
