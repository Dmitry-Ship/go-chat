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

const (
	ParticipantTypeJoined  = "joined"
	ParticipantTypePrivate = "private"
	ParticipantTypeOwner   = "owner"
)

type Participant struct {
	aggregate
	ID             uuid.UUID
	ConversationID uuid.UUID
	UserID         uuid.UUID
	CreatedAt      time.Time
	Type           string
	IsActive       bool
}

type ParticipantLeaver interface {
	LeavePublicConversation(conversationID uuid.UUID) error
}

func newParticipant(conversationId uuid.UUID, userId uuid.UUID, participantType string) *Participant {
	return &Participant{
		ID:             uuid.New(),
		ConversationID: conversationId,
		UserID:         userId,
		CreatedAt:      time.Now(),
		Type:           participantType,
		IsActive:       true,
	}
}

func NewOwnerParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	participant := newParticipant(conversationId, userId, ParticipantTypeOwner)

	return participant
}

func NewJoinedParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	participant := newParticipant(conversationId, userId, ParticipantTypeJoined)

	participant.AddEvent(NewPublicConversationJoined(conversationId, userId))

	return participant
}

func NewPrivateParticipant(conversationId uuid.UUID, userId uuid.UUID) *Participant {
	participant := newParticipant(conversationId, userId, ParticipantTypePrivate)
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
