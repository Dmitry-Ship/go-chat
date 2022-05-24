package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ParticipantRepository interface {
	Store(participant *Participant) error
	GetByConversationIDAndUserID(conversationID uuid.UUID, userID uuid.UUID) (*Participant, error)
	Update(participant *Participant) error
}

type Participant struct {
	aggregate
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	CreatedAt      time.Time
	IsActive       bool
}

type ParticipantLeaver interface {
	LeaveGroupConversation(conversationID uuid.UUID) error
}

func NewParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	return &Participant{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userId,
		CreatedAt:      time.Now(),
		IsActive:       true,
	}
}

func (participant *Participant) LeaveGroupConversation(conversationID uuid.UUID) error {
	if participant.ConversationID != conversationID {
		return errors.New("participant is not in conversation")
	}

	if !participant.IsActive {
		return errors.New("participant is already left")
	}

	participant.AddEvent(NewGroupConversationLeft(conversationID, participant.UserID))

	return nil
}
