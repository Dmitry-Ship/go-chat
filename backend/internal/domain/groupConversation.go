package domain

import (
	"errors"

	"github.com/google/uuid"
)

type GroupConversationRepository interface {
	Store(conversation *GroupConversation) error
	Update(conversation *GroupConversation) error
	GetByID(id uuid.UUID) (*GroupConversation, error)
}

type GroupConversation struct {
	Conversation
	ID     uuid.UUID
	Name   string
	Avatar string
	Owner  Participant
}

func NewGroupConversation(id uuid.UUID, name string, creatorID uuid.UUID) (*GroupConversation, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}

	if len(name) > 100 {
		return nil, errors.New("name is too long")
	}

	groupConversation := &GroupConversation{
		Conversation: Conversation{
			ID:       id,
			Type:     "group",
			IsActive: true,
		},
		ID:     uuid.New(),
		Name:   name,
		Avatar: string(name[0]),
		Owner:  *NewParticipant(id, creatorID),
	}

	groupConversation.AddEvent(newGroupConversationCreatedEvent(id, creatorID))

	return groupConversation, nil
}

func (groupConversation *GroupConversation) Delete(userID uuid.UUID) error {
	if groupConversation.Owner.UserID != userID {
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
	if newName == "" {
		return errors.New("name is empty")
	}

	if len(newName) > 100 {
		return errors.New("name is too long")
	}

	if groupConversation.Owner.UserID != userID {
		return errors.New("user is not owner")
	}

	groupConversation.Name = newName
	groupConversation.Avatar = string(newName[0])

	groupConversation.AddEvent(newGroupConversationRenamedEvent(groupConversation.Conversation.ID, userID, newName))

	return nil
}

func (groupConversation *GroupConversation) Join(userID uuid.UUID) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	participant := NewParticipant(groupConversation.Conversation.ID, userID)

	participant.AddEvent(newGroupConversationJoinedEvent(groupConversation.Conversation.ID, userID))

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

	if groupConversation.Owner.UserID == inviteeID {
		return nil, errors.New("user is owner")
	}

	if userID == inviteeID {
		return nil, errors.New("cannot invite yourself")
	}

	participant := NewParticipant(groupConversation.Conversation.ID, inviteeID)

	participant.AddEvent(newGroupConversationInvitedEvent(groupConversation.Conversation.ID, userID, inviteeID))

	return participant, nil
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

func (groupConversation *GroupConversation) SendJoinedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newJoinedConversationMessage(conversationID, userID)

	return message, nil
}

func (groupConversation *GroupConversation) SendInvitedConversationMessage(conversationID uuid.UUID, userID uuid.UUID) (*Message, error) {
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

func (groupConversation *GroupConversation) SendLeftConversationMessage(conversationID uuid.UUID, userID uuid.UUID) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newLeftConversationMessage(conversationID, userID)

	return message, nil
}
