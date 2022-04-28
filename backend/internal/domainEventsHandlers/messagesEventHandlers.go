package domainEventsHandlers

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/services"
	"log"
)

type messagesEventHandlers struct {
	pubsub               domain.EventsSubscriber
	messagingService     services.MessagingService
	notificationsService services.NotificationService
}

func NewMessagesEventHandlers(pubsub domain.EventsSubscriber, messagingService services.MessagingService) *messagesEventHandlers {
	return &messagesEventHandlers{
		pubsub:           pubsub,
		messagingService: messagingService,
	}
}

func (h *messagesEventHandlers) Run() {
	for {
		select {
		case event := <-h.pubsub.Subscribe(domain.PublicConversationRenamedName):
			e, ok := event.(*domain.PublicConversationRenamed)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.messagingService.SendRenamedConversationMessage(e.ConversationID, e.UserID, e.NewName)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

		case event := <-h.pubsub.Subscribe(domain.PublicConversationLeftName):
			e, ok := event.(*domain.PublicConversationLeft)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.messagingService.SendLeftConversationMessage(e.ConversationID, e.UserID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

		case event := <-h.pubsub.Subscribe(domain.PublicConversationJoinedName):
			e, ok := event.(*domain.PublicConversationJoined)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.messagingService.SendJoinedConversationMessage(e.ConversationID, e.UserID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}
		}

	}
}
