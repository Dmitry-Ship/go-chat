package services

import (
	"context"
	"errors"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

var ErrSystemMessageValidationFailed = errors.New("system message validation failed: conversation or user not found")

type SystemMessageService interface {
	SaveJoinedMessage(ctx context.Context, conversationID, userID uuid.UUID) error
	SaveLeftMessage(ctx context.Context, conversationID, userID uuid.UUID) error
	SaveInvitedMessage(ctx context.Context, conversationID, userID uuid.UUID) error
	SaveRenamedMessage(ctx context.Context, conversationID, userID uuid.UUID, newName string) error
}

type systemMessageService struct {
	messages domain.MessageRepository
}

func NewSystemMessageService(
	messages domain.MessageRepository,
) *systemMessageService {
	return &systemMessageService{
		messages: messages,
	}
}

func (s *systemMessageService) SaveJoinedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeJoinedConversation, "")

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store joined message error: %w", err)
	}

	if !stored {
		return ErrSystemMessageValidationFailed
	}

	return nil
}

func (s *systemMessageService) SaveInvitedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeInvitedConversation, "")

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store invited message error: %w", err)
	}

	if !stored {
		return ErrSystemMessageValidationFailed
	}

	return nil
}

func (s *systemMessageService) SaveRenamedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, newName string) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeRenamedConversation, newName)

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store renamed message error: %w", err)
	}

	if !stored {
		return ErrSystemMessageValidationFailed
	}

	return nil
}

func (s *systemMessageService) SaveLeftMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	message := domain.NewSystemMessage(uuid.New(), conversationID, userID, domain.MessageTypeLeftConversation, "")

	stored, err := s.messages.StoreSystemMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("store left message error: %w", err)
	}

	if !stored {
		return ErrSystemMessageValidationFailed
	}

	return nil
}
