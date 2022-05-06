package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"
	"log"

	"github.com/google/uuid"
)

type eventHandlers struct {
	pubsub   domain.EventsSubscriber
	commands *app.Commands
	queries  readModel.QueriesRepository
}

func NewEventHandlers(pubsub domain.EventsSubscriber, commands *app.Commands, queries readModel.QueriesRepository) *eventHandlers {
	return &eventHandlers{
		pubsub:   pubsub,
		commands: commands,
		queries:  queries,
	}
}

func (h *eventHandlers) ListenForEvents() {
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
}

func (h *eventHandlers) HandlePublicConversationLeft() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationLeftName) {
		e, ok := event.(*domain.PublicConversationLeft)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
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
}

func (h *eventHandlers) HandlePublicConversationJoined() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationJoinedName) {
		e, ok := event.(*domain.PublicConversationJoined)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
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

		err := h.commands.NotificationTopicService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}

		err = h.commands.NotificationTopicService.DeleteTopic("conversation:" + e.ConversationID.String())

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

		messageDTO, err := h.queries.GetNotificationMessage(e.MessageID, e.UserID)

		if err != nil {
			log.Println(err)
			continue
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
}

func (h *eventHandlers) HandlePublicConversationCreated() {
	for event := range h.pubsub.Subscribe(domain.PublicConversationCreatedName) {
		e, ok := event.(*domain.PublicConversationCreated)

		if !ok {
			log.Printf("Wrong event type: %T", e)
			continue
		}

		err := h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.OwnerID)

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

		err := h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.FromUserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}

		err = h.commands.NotificationTopicService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.ToUserID)

		if err != nil {
			log.Printf("Error handling %s event: %v", e.GetName(), err)
		}
	}
}
