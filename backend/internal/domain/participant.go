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
	LeavePublicConversation(conversationID uuid.UUID) error
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

func NewJoinedParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	participant := NewParticipant(conversationId, userId)

	participant.AddEvent(NewPublicConversationJoined(conversationId, userId))

	return participant
}

func NewInvitedParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	participant := NewParticipant(conversationId, userId)

	participant.AddEvent(NewPublicConversationInvited(conversationId, userId))

	return participant
}

func (participant *Participant) LeavePublicConversation(conversationID uuid.UUID) error {
	if participant.ConversationID != conversationID {
		return errors.New("participant is not in conversation")
	}

	if !participant.IsActive {
		return errors.New("participant is already left")
	}

	participant.AddEvent(NewPublicConversationLeft(conversationID, participant.UserID))

	return nil
}
