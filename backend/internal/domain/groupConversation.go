package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type GroupConversationRepository interface {
	Store(conversation *GroupConversation) error
	Update(conversation *GroupConversation) error
	GetByID(id uuid.UUID) (*GroupConversation, error)
}

type GroupConversationData struct {
	ID     uuid.UUID
	Name   string
	Avatar string
	Owner  Participant
}
type GroupConversation struct {
	Conversation
	Data GroupConversationData
}

func NewGroupConversation(id uuid.UUID, name string, creatorID uuid.UUID) (*GroupConversation, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}

	groupConversation := &GroupConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "group",
			CreatedAt: time.Now(),
			IsActive:  true,
		},
		Data: GroupConversationData{
			ID:     uuid.New(),
			Name:   name,
			Avatar: string(name[0]),
			Owner:  *NewParticipant(id, creatorID),
		},
	}

	groupConversation.AddEvent(newGroupConversationCreatedEvent(id, creatorID))

	return groupConversation, nil
}

func (groupConversation *GroupConversation) Delete(userID uuid.UUID) error {
	if groupConversation.Data.Owner.UserID != userID {
		return errors.New("user is not owner")
	}

	if !groupConversation.IsActive {
		return errors.New("conversation is not active")
	}

	groupConversation.IsActive = false

	groupConversation.AddEvent(newGroupConversationDeletedEvent(groupConversation.Conversation.ID))

	return nil
}

func (groupConversation *GroupConversation) Rename(newName string, userID uuid.UUID) error {
	if groupConversation.Data.Owner.UserID != userID {
		return errors.New("user is not owner")
	}

	groupConversation.Data.Name = newName
	groupConversation.Data.Avatar = string(newName[0])

	groupConversation.AddEvent(newGroupConversationRenamedEvent(groupConversation.ID, userID, newName))

	return nil
}

func (groupConversation *GroupConversation) SendTextMessage(text string, participant *Participant) (*TextMessage, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if participant.ConversationID != groupConversation.Conversation.ID {
		return nil, errors.New("user is not participant")
	}

	if !participant.IsActive {
		return nil, errors.New("user is not participant")
	}

	message, err := newTextMessage(groupConversation.Conversation.ID, participant.UserID, text)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (groupConversation *GroupConversation) Join(userID uuid.UUID) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	participant := NewParticipant(groupConversation.ID, userID)

	participant.AddEvent(newGroupConversationJoinedEvent(groupConversation.ID, userID))

	return participant, nil
}

func (groupConversation *GroupConversation) Leave(participant *Participant) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if participant.ConversationID != groupConversation.Conversation.ID {
		return nil, errors.New("participant is not in conversation")
	}

	if !participant.IsActive {
		return nil, errors.New("participant already left")
	}

	participant.IsActive = false

	participant.AddEvent(newGroupConversationLeftEvent(groupConversation.Conversation.ID, participant.UserID))

	return participant, nil
}

func (groupConversation *GroupConversation) Invite(userID uuid.UUID, inviteeID uuid.UUID) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if groupConversation.Data.Owner.UserID == inviteeID {
		return nil, errors.New("user is owner")
	}

	if userID == inviteeID {
		return nil, errors.New("cannot invite yourself")
	}

	participant := NewParticipant(groupConversation.ID, inviteeID)

	participant.AddEvent(newGroupConversationInvitedEvent(groupConversation.ID, userID, inviteeID))

	return participant, nil
}

func (groupConversation *GroupConversation) SendJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) (*JoinedConversationMessage, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newJoinedConversationMessage(conversationID, userID)

	return message, nil
}

func (groupConversation *GroupConversation) SendInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) (*InvitedConversationMessage, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newInvitedConversationMessage(conversationID, userID)

	return message, nil
}

func (groupConversation *GroupConversation) SendRenamedConversationMessage(conversationID uuid.UUID, userID uuid.UUID, newName string) (*ConversationRenamedMessage, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newConversationRenamedMessage(conversationID, userID, newName)

	return message, nil
}

func (groupConversation *GroupConversation) SendLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) (*LeftConversationMessage, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newLeftConversationMessage(conversationID, userID)

	return message, nil
}
