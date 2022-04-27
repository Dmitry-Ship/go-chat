package services

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"GitHub/go-chat/backend/internal/readModel"
	"log"

	"github.com/google/uuid"
)

type notificationEventHandlers struct {
	pubsub               domain.PubSub
	messages             readModel.MessageQueryRepository
	notificationsService NotificationService
}

func NewNotificationEventHandlers(pubsub domain.PubSub, notificationsService NotificationService, messages readModel.MessageQueryRepository) *notificationEventHandlers {
	return &notificationEventHandlers{
		pubsub:               pubsub,
		messages:             messages,
		notificationsService: notificationsService,
	}
}

func (h *notificationEventHandlers) Run() {
	for {
		select {
		case event := <-h.pubsub.Subscribe(domain.PublicConversationRenamedName):
			e, ok := event.(*domain.PublicConversationRenamed)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
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
			h.notificationsService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

		case event := <-h.pubsub.Subscribe(domain.PublicConversationDeletedName):
			e, ok := event.(*domain.PublicConversationDeleted)

			if !ok {
				log.Printf("Unknown event type: %T", e)
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

			h.notificationsService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

			err := h.notificationsService.DeleteTopic("conversation:" + e.ConversationID.String())

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

		case event := <-h.pubsub.Subscribe(domain.MessageSentName):
			e, ok := event.(*domain.MessageSent)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			messageDTO, err := h.messages.GetMessageByID(e.MessageID, e.UserID)

			if err != nil {
				log.Println(err)
				continue
			}

			notification := ws.OutgoingNotification{
				Type:    "message",
				Payload: messageDTO,
			}

			h.notificationsService.BroadcastToTopic("conversation:"+e.ConversationID.String(), notification)

		case event := <-h.pubsub.Subscribe(domain.PublicConversationCreatedName):
			e, ok := event.(*domain.PublicConversationCreated)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.notificationsService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.OwnerID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

		case event := <-h.pubsub.Subscribe(domain.PublicConversationJoinedName):
			e, ok := event.(*domain.PublicConversationJoined)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.notificationsService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.UserID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

		case event := <-h.pubsub.Subscribe(domain.PublicConversationLeftName):
			e, ok := event.(*domain.PublicConversationLeft)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.notificationsService.UnsubscribeFromTopic("conversation:"+e.ConversationID.String(), e.UserID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

		case event := <-h.pubsub.Subscribe(domain.PrivateConversationCreatedName):
			e, ok := event.(*domain.PrivateConversationCreated)

			if !ok {
				log.Printf("Unknown event type: %T", e)
				continue
			}

			err := h.notificationsService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.FromUserID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}

			err = h.notificationsService.SubscribeToTopic("conversation:"+e.ConversationID.String(), e.ToUserID)

			if err != nil {
				log.Printf("Error handling %s event: %v", e.GetName(), err)
			}
		}
	}
}
