package services

import (
	"GitHub/go-chat/backend/internal/domain"
	ws "GitHub/go-chat/backend/internal/websocket"
	"context"

	"github.com/google/uuid"
)

type NotificationsPipeline interface {
	GetReceivers(ctx context.Context, event domain.DomainEvent) (chan uuid.UUID, error)
	BuildMessage(ctx context.Context, receiverIDsChan <-chan uuid.UUID, event domain.DomainEvent) (chan ws.OutgoingNotification, chan error)
	BroadcastMessage(ctx context.Context, messageChan <-chan ws.OutgoingNotification) chan error
}

type notificationsPipeline struct {
	notificationCommands NotificationService
	notificationResolver NotificationResolverService
	notificationBuilder  NotificationBuilderService
}

func NewNotificationsPipeline(
	notificationCommands NotificationService,
	notificationResolver NotificationResolverService,
	notificationBuilder NotificationBuilderService,
) *notificationsPipeline {
	return &notificationsPipeline{
		notificationCommands: notificationCommands,
		notificationResolver: notificationResolver,
		notificationBuilder:  notificationBuilder,
	}
}

func (h *notificationsPipeline) GetReceivers(ctx context.Context, event domain.DomainEvent) (chan uuid.UUID, error) {
	receiversIDs, err := h.notificationResolver.GetReceiversFromEvent(event)
	receiverIDsChan := make(chan uuid.UUID, len(receiversIDs))
	defer close(receiverIDsChan)

	if err != nil {
		return nil, err
	}

	for _, receiverID := range receiversIDs {
		select {
		case receiverIDsChan <- receiverID:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return receiverIDsChan, nil
}

func (h *notificationsPipeline) BuildMessage(ctx context.Context, receiverIDsChan <-chan uuid.UUID, event domain.DomainEvent) (chan ws.OutgoingNotification, chan error) {
	messageChan := make(chan ws.OutgoingNotification, len(receiverIDsChan))
	errorChan := make(chan error, len(receiverIDsChan))

	go func() {
		defer close(messageChan)
		defer close(errorChan)
		for receiverID := range receiverIDsChan {
			message, err := h.notificationBuilder.BuildMessageFromEvent(receiverID, event)

			if err != nil {
				errorChan <- err
				return
			}

			select {
			case messageChan <- message:
			case <-ctx.Done():
				return

			}
		}
	}()

	return messageChan, errorChan
}

func (h *notificationsPipeline) BroadcastMessage(ctx context.Context, messageChan <-chan ws.OutgoingNotification) chan error {
	errorChan := make(chan error, len(messageChan))
	defer close(errorChan)

	for message := range messageChan {
		select {
		case <-ctx.Done():
			return errorChan
		default:
			err := h.notificationCommands.Send(message)
			if err != nil {
				errorChan <- err
				return errorChan
			}
		}
	}

	return errorChan
}
