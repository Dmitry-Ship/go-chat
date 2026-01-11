package services

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
)

type MessageService interface {
	SendTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error
}

type messageService struct {
	queries       readModel.QueriesRepository
	notifications NotificationService
}

func NewMessageService(
	queries readModel.QueriesRepository,
	notifications NotificationService,
) MessageService {
	return &messageService{
		queries:       queries,
		notifications: notifications,
	}
}

func (s *messageService) SendTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error {
	isMember, err := s.queries.IsMember(conversationID, userID)
	if err != nil {
		return fmt.Errorf("is member error: %w", err)
	}
	if !isMember {
		return fmt.Errorf("user is not in conversation: %w", domain.ErrorUserNotInConversation)
	}

	messageID := uuid.New()

	if _, err := domain.NewTextMessageContent(messageText); err != nil {
		return fmt.Errorf("validate message text error: %w", err)
	}

	messageDTO, err := s.queries.StoreMessageAndReturnWithUser(messageID, conversationID, userID, messageText, 0)
	if err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	if err := s.notifications.Broadcast(ctx, conversationID, ws.OutgoingNotification{Type: "message", Payload: messageDTO, UserID: userID}); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}

	return nil
}
