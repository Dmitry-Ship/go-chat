package domain

import (
	"errors"

	"github.com/google/uuid"
)

type DirectConversationRepository interface {
	Store(conversation *DirectConversation) error
	GetID(firstUserID uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*DirectConversation, error)
}

type DirectConversation struct {
	Conversation
	ID       uuid.UUID
	ToUser   Participant
	FromUser Participant
}

func NewDirectConversation(id uuid.UUID, to uuid.UUID, from uuid.UUID) (*DirectConversation, error) {
	if to == from {
		return nil, errors.New("cannot chat with yourself")
	}

	directConversation := DirectConversation{
		Conversation: Conversation{
			ID:       id,
			Type:     ConversationTypeDirect,
			IsActive: true,
		},
		ID:       uuid.New(),
		ToUser:   *NewParticipant(id, to),
		FromUser: *NewParticipant(id, from),
	}

	directConversation.AddEvent(newDirectConversationCreatedEvent(id, to, from))

	return &directConversation, nil
}

func (directConversation *DirectConversation) GetFromUser() *Participant {
	return &directConversation.FromUser
}

func (directConversation *DirectConversation) GetToUser() *Participant {
	return &directConversation.ToUser
}

func (directConversation *DirectConversation) SendTextMessage(text string, userID uuid.UUID) (*TextMessage, error) {
	if directConversation.ToUser.UserID != userID && directConversation.FromUser.UserID != userID {
		return nil, errors.New("user is not participant")
	}

	message, err := newTextMessage(directConversation.ID, userID, text)

	if err != nil {
		return nil, err
	}

	return message, nil
}