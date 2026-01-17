package services

import (
	"context"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/repository"
	ws "GitHub/go-chat/backend/internal/websocket"
)

type MessageService interface {
	Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error)
}

type messageService struct {
	messages      repository.MessageRepository
	notifications NotificationService
}

func NewMessageService(
	messages repository.MessageRepository,
	notifications NotificationService,
) MessageService {
	return &messageService{
		messages:      messages,
		notifications: notifications,
	}
}

func (s *messageService) Send(ctx context.Context, message *domain.Message) (readModel.MessageDTO, error) {
	dto, err := s.messages.Send(ctx, message)
	if err != nil {
		return dto, err
	}

	if err := s.notifications.Broadcast(ctx, message.ConversationID, ws.OutgoingNotification{Type: "message", Payload: dto, UserID: message.UserID}); err != nil {
		return dto, err
	}

	return dto, nil
}
