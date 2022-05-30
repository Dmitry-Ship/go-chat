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
	Participants []Participant
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
		Participants: []Participant{
			*NewParticipant(uuid.New(), id, to),
			*NewParticipant(uuid.New(), id, from),
		},
	}

	directConversation.AddEvent(newDirectConversationCreatedEvent(id, []uuid.UUID{to, from}))

	return &directConversation, nil
}

func (directConversation *DirectConversation) SendTextMessage(messageID uuid.UUID, text string, userID uuid.UUID) (*Message, error) {
	isParticipant := false
	for _, participant := range directConversation.Participants {
		if participant.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("user is not participant")
	}

	content, err := newTextMessageContent(text)

	if err != nil {
		return nil, err
	}

	message := newTextMessage(messageID, directConversation.Conversation.ID, userID, content)

	return message, nil
}
