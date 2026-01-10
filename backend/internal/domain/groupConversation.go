package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type GroupConversationRepository interface {
	Store(ctx context.Context, conversation *GroupConversation) error
	Update(ctx context.Context, conversation *GroupConversation) error
	Rename(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*GroupConversation, error)
}

var (
	ErrorUserNotInConversation = errors.New("user is not in conversation")
	ErrorUserNotOwner          = errors.New("user is not owner")
	ErrorCannotInviteOneself   = errors.New("cannot invite yourself")
	ErrorCannotKickOneself     = errors.New("cannot kick yourself")
	ErrorOwnerCannotLeave      = errors.New("owner cannot leave conversation")
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

func NewGroupConversation(id uuid.UUID, name string, creatorId uuid.UUID) (*GroupConversation, error) {
	groupConversation := &GroupConversation{
		Conversation: Conversation{
			ID:   id,
			Type: ConversationTypeGroup,
		},
		ID:     uuid.New(),
		Name:   name,
		Avatar: string(name[0]),
		Owner:  *NewParticipant(uuid.New(), id, creatorId),
	}

	return groupConversation, nil
}

func (groupConversation *GroupConversation) isJoined(participant *Participant) bool {
	return participant.ConversationID == groupConversation.Conversation.ID
}

func (groupConversation *GroupConversation) Delete(participant *Participant) error {
	if !groupConversation.isJoined(participant) {
		return ErrorUserNotInConversation
	}

	if groupConversation.Owner.UserID != participant.UserID {
		return ErrorUserNotOwner
	}

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

	return nil
}

func (groupConversation *GroupConversation) Invite(inviter *Participant, invitee *User) (*Participant, error) {
	if !groupConversation.isJoined(inviter) {
		return nil, ErrorUserNotInConversation
	}

	if inviter.UserID == invitee.ID {
		return nil, ErrorCannotInviteOneself
	}

	participant := NewParticipant(uuid.New(), groupConversation.Conversation.ID, invitee.ID)

	return participant, nil
}

func (groupConversation *GroupConversation) Kick(kicker *Participant, target *Participant) (*Participant, error) {
	if !groupConversation.isJoined(target) {
		return nil, ErrorUserNotInConversation
	}

	if groupConversation.Owner.UserID != kicker.UserID {
		return nil, ErrorUserNotOwner
	}

	if kicker.UserID == target.UserID {
		return nil, ErrorCannotKickOneself
	}

	return target, nil
}

func (groupConversation *GroupConversation) SendJoinedConversationMessage(messageID uuid.UUID, user *User) (*Message, error) {
	message := newJoinedConversationMessage(messageID, groupConversation.Conversation.ID, user.ID)

	return message, nil
}

func (groupConversation *GroupConversation) SendInvitedConversationMessage(messageID uuid.UUID, user *User) (*Message, error) {
	message := newInvitedConversationMessage(messageID, groupConversation.Conversation.ID, user.ID)

	return message, nil
}

func (groupConversation *GroupConversation) SendRenamedConversationMessage(messageID uuid.UUID, participant *Participant, newName string) (*Message, error) {
	content := newRenamedMessageContent(newName)
	message := newConversationRenamedMessage(messageID, groupConversation.Conversation.ID, participant.UserID, content)

	return message, nil
}

func (groupConversation *GroupConversation) SendLeftConversationMessage(messageID uuid.UUID, participant *Participant) (*Message, error) {
	message := newLeftConversationMessage(messageID, groupConversation.Conversation.ID, participant.UserID)

	return message, nil
}
