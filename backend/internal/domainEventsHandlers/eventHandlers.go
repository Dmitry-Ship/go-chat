package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"log"

	"github.com/google/uuid"
)

type eventHandlers struct {
	pubsub   domain.EventsSubscriber
	Commands *app.Commands
	Queries  *app.Queries
}

func NewEventHandlers(pubsub domain.EventsSubscriber, commands *app.Commands, queries *app.Queries) *eventHandlers {
	return &eventHandlers{
		pubsub:   pubsub,
		Commands: commands,
		Queries:  queries,
	}
}

func (h *eventHandlers) ListerForEvents() {
	go h.HandleMessageSent()
	go h.HandlePublicConversationCreated()
	go h.HandlePublicConversationRenamed()
	go h.HandlePublicConversationLeft()
	go h.HandlePublicConversationJoined()
	go h.HandlePublicConversationDeleted()
	go h.HandlePrivateConversationCreated()
}

func (h *eventHandlers) HandlePublicConversationRenamed() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationRenamedName) {
		e, ok := event.(*domain.PublicConversationRenamed)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		err := h.Commands.MessagingService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

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
		err = h.Commands.NotificationService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}

func (h *eventHandlers) HandlePublicConversationLeft() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationLeftName) {
		e, ok := event.(*domain.PublicConversationLeft)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		err := h.Commands.MessagingService.SendLeftConversationMessage(e.ConversationID, e.UserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}

		err = h.Commands.NotificationService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}

func (h *eventHandlers) HandlePublicConversationJoined() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationJoinedName) {
		e, ok := event.(*domain.PublicConversationJoined)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		err := h.Commands.MessagingService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}

		err = h.Commands.NotificationService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.UserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}

func (h *eventHandlers) HandlePublicConversationDeleted() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationDeletedName) {
		e, ok := event.(*domain.PublicConversationDeleted)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		notification := ws.OutgoingNotification{
			Type: "conversation_deleted",
			Payload: struct {
				ConversationId uuid.UUID `json:"conversation_id"`
			}{
				ConversationId: e.ConversationID,
			},
		}

		err := h.Commands.NotificationService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}

		err = h.Commands.NotificationService.DeleteTopic("conversation:" + e.ConversationID.String())

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}

func (h *eventHandlers) HandleMessageSent() {
	for event := range h.pubsub.Subscribe(domain.MessageSentName) {
		e, ok := event.(*domain.MessageSent)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		messageDTO, err := h.Queries.Messages.GetNotificationMessage(e.MessageID, e.UserID)

		if err != nil {
			log.Println(err)
			continue
		}

		notification := ws.OutgoingNotification{
			Type:    "message",
			Payload: messageDTO,
		}

		err = h.Commands.NotificationService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}

func (h *eventHandlers) HandlePublicConversationCreated() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationCreatedName) {
		e, ok := event.(*domain.PublicConversationCreated)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		err := h.Commands.NotificationService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.OwnerID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}

func (h *eventHandlers) HandlePrivateConversationCreated() {
	for event := range h.pubsub.Subscribe(domain.PrivateConversationCreatedName) {
		e, ok := event.(*domain.PrivateConversationCreated)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		err := h.Commands.NotificationService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.FromUserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}

		err = h.Commands.NotificationService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.ToUserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}
