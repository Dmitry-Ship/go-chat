package domain

import (
	"errors"

	"github.com/google/uuid"
)

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
			ID:   id,
			Type: ConversationTypeDirect,
		},
		Participants: []Participant{
			*NewParticipant(uuid.New(), id, to),
			*NewParticipant(uuid.New(), id, from),
		},
	}

	return &directConversation, nil
}
