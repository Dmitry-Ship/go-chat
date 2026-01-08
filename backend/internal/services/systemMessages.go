package services

import (
	"context"
	"fmt"

	"GitHub/go-chat/backend/internal/domain"

	"github.com/google/uuid"
)

type SystemMessageService interface {
	SaveJoinedMessage(ctx context.Context, conversationID, userID uuid.UUID) error
	SaveLeftMessage(ctx context.Context, conversationID, userID uuid.UUID) error
	SaveInvitedMessage(ctx context.Context, conversationID, userID uuid.UUID) error
	SaveRenamedMessage(ctx context.Context, conversationID, userID uuid.UUID, newName string) error
}

type systemMessageService struct {
	groupConversations domain.GroupConversationRepository
	users              domain.UserRepository
	participants       domain.ParticipantRepository
	messages           domain.MessageRepository
}

func NewSystemMessageService(
	groupConversations domain.GroupConversationRepository,
	users domain.UserRepository,
	participants domain.ParticipantRepository,
	messages domain.MessageRepository,
) *systemMessageService {
	return &systemMessageService{
		groupConversations: groupConversations,
		users:              users,
		participants:       participants,
		messages:           messages,
	}
}

func (s *systemMessageService) SaveJoinedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	message, err := conversation.SendJoinedConversationMessage(messageID, user)
	if err != nil {
		return fmt.Errorf("send joined message error: %w", err)
	}

	if err := s.messages.Store(ctx, message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *systemMessageService) SaveInvitedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id error: %w", err)
	}

	message, err := conversation.SendInvitedConversationMessage(messageID, user)
	if err != nil {
		return fmt.Errorf("send invited message error: %w", err)
	}

	if err := s.messages.Store(ctx, message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *systemMessageService) SaveRenamedMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, newName string) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)
	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	message, err := conversation.SendRenamedConversationMessage(messageID, participant, newName)
	if err != nil {
		return fmt.Errorf("send renamed message error: %w", err)
	}

	if err := s.messages.Store(ctx, message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}

func (s *systemMessageService) SaveLeftMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	conversation, err := s.groupConversations.GetByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation by id error: %w", err)
	}

	messageID := uuid.New()

	participant, err := s.participants.GetByConversationIDAndUserID(ctx, conversationID, userID)
	if err != nil {
		return fmt.Errorf("get participant error: %w", err)
	}

	message, err := conversation.SendLeftConversationMessage(messageID, participant)
	if err != nil {
		return fmt.Errorf("send left message error: %w", err)
	}

	if err := s.messages.Store(ctx, message); err != nil {
		return fmt.Errorf("store message error: %w", err)
	}

	return nil
}
