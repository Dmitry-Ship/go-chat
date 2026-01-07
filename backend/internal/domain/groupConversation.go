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

var (
	ErrorUserNotInConversation = errors.New("user is not in conversation")
	ErrorConversationNotActive = errors.New("conversation is not active")
	ErrorUserNotOwner          = errors.New("user is not owner")
	ErrorCannotInviteOneself   = errors.New("cannot invite yourself")
	ErrorCannotKickOneself     = errors.New("cannot kick yourself")
)

func ValidateConversationName(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}

	if len(name) > 100 {
		return errors.New("name is too long")
	}

	return nil
}

type GroupConversation struct {
	Conversation
	ID     uuid.UUID
	Name   string
	Avatar string
	Owner  Participant
}

func NewGroupConversation(id uuid.UUID, name string, creator User) (*GroupConversation, error) {
	groupConversation := &GroupConversation{
		Conversation: Conversation{
			ID:       id,
			Type:     ConversationTypeGroup,
			IsActive: true,
		},
		ID:     uuid.New(),
		Name:   name,
		Avatar: string(name[0]),
		Owner:  *NewParticipant(uuid.New(), id, creator.ID),
	}

	groupConversation.AddEvent(newGroupConversationCreatedEvent(id, creator.ID))

	return groupConversation, nil
}

func (groupConversation *GroupConversation) isJoined(participant *Participant) bool {
	return participant.ConversationID == groupConversation.Conversation.ID && participant.IsActive
}

func (groupConversation *GroupConversation) Delete(participant *Participant) error {
	if !groupConversation.isJoined(participant) {
		return ErrorUserNotInConversation
	}

	if groupConversation.Owner.UserID != participant.UserID {
		return ErrorUserNotOwner
	}

	if !groupConversation.IsActive {
		return ErrorConversationNotActive
	}

	groupConversation.IsActive = false

	groupConversation.AddEvent(newGroupConversationDeletedEvent(groupConversation.Conversation.ID))

	return nil
}

func (groupConversation *GroupConversation) Rename(newName string, participant *Participant) error {
	if !groupConversation.isJoined(participant) {
		return ErrorUserNotInConversation
	}

	if groupConversation.Owner.UserID != participant.UserID {
		return ErrorUserNotOwner
	}

	groupConversation.Name = newName
	groupConversation.Avatar = string(newName[0])

	groupConversation.AddEvent(newGroupConversationRenamedEvent(groupConversation.Conversation.ID, participant.UserID, newName))

	return nil
}

func (groupConversation *GroupConversation) Join(user User) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	participant := NewParticipant(uuid.New(), groupConversation.Conversation.ID, user.ID)

	participant.AddEvent(newGroupConversationJoinedEvent(groupConversation.Conversation.ID, user.ID))

	return participant, nil
}

func (groupConversation *GroupConversation) Leave(participant *Participant) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	if !groupConversation.isJoined(participant) {
		return nil, ErrorUserNotInConversation
	}

	participant.IsActive = false

	participant.AddEvent(newGroupConversationLeftEvent(groupConversation.Conversation.ID, participant.UserID))

	return participant, nil
}

func (groupConversation *GroupConversation) Invite(inviter *Participant, invitee *User) (*Participant, error) {
	if !groupConversation.isJoined(inviter) {
		return nil, ErrorUserNotInConversation
	}

	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	if inviter.UserID == invitee.ID {
		return nil, ErrorCannotInviteOneself
	}

	participant := NewParticipant(uuid.New(), groupConversation.Conversation.ID, invitee.ID)

	participant.AddEvent(newGroupConversationInvitedEvent(groupConversation.Conversation.ID, inviter.UserID, invitee.ID))

	return participant, nil
}

func (groupConversation *GroupConversation) Kick(kicker *Participant, target *Participant) (*Participant, error) {
	if !groupConversation.isJoined(target) {
		return nil, ErrorUserNotInConversation
	}

	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	if groupConversation.Owner.UserID != kicker.UserID {
		return nil, ErrorUserNotOwner
	}

	if kicker.UserID == target.UserID {
		return nil, ErrorCannotKickOneself
	}

	target.IsActive = false

	target.AddEvent(newGroupConversationLeftEvent(groupConversation.Conversation.ID, target.UserID))

	return target, nil
}

func (groupConversation *GroupConversation) SendTextMessage(messageID uuid.UUID, text string, participant *Participant) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	if !groupConversation.isJoined(participant) {
		return nil, ErrorUserNotInConversation
	}

	content, err := newTextMessageContent(text)

	if err != nil {
		return nil, err
	}

	message := newTextMessage(messageID, groupConversation.Conversation.ID, participant.UserID, content)

	return message, nil
}

func (groupConversation *GroupConversation) SendJoinedConversationMessage(messageID uuid.UUID, user *User) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	message := newJoinedConversationMessage(messageID, groupConversation.Conversation.ID, user.ID)

	return message, nil
}

func (groupConversation *GroupConversation) SendInvitedConversationMessage(messageID uuid.UUID, user *User) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	message := newInvitedConversationMessage(messageID, groupConversation.Conversation.ID, user.ID)

	return message, nil
}

func (groupConversation *GroupConversation) SendRenamedConversationMessage(messageID uuid.UUID, participant *Participant, newName string) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, ErrorConversationNotActive
	}

	content := newRenamedMessageContent(newName)
	message := newConversationRenamedMessage(messageID, groupConversation.Conversation.ID, participant.UserID, content)

	return message, nil
}

func (groupConversation *GroupConversation) SendLeftConversationMessage(messageID uuid.UUID, participant *Participant) (*Message, error) {
	message := newLeftConversationMessage(messageID, groupConversation.Conversation.ID, participant.UserID)

	return message, nil
}
