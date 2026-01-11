package services

import (
	"context"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type MessageService interface {
	Send(ctx context.Context, message *domain.Message, requestUserID uuid.UUID) (readModel.MessageDTO, error)
}

type messageService struct {
	messages      domain.MessageRepository
	notifications NotificationService
}

func NewMessageService(
	messages domain.MessageRepository,
	notifications NotificationService,
) MessageService {
	return &messageService{
		messages:      messages,
		notifications: notifications,
	}
}

func (s *messageService) Send(ctx context.Context, message *domain.Message, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	dto, err := s.messages.Send(ctx, message, requestUserID)
	if err != nil {
		return dto, err
	}

	if err := s.notifications.Broadcast(ctx, message.ConversationID, ws.OutgoingNotification{Type: "message", Payload: dto, UserID: message.UserID}); err != nil {
		return dto, err
	}

	return dto, nil
}
