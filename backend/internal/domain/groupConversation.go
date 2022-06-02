package domain

import (
	"errors"

	"github.com/google/uuid"
)

type GroupConversationRepository interface {
	GenericRepository[*GroupConversation]
	GetByID(id uuid.UUID) (*GroupConversation, error)
}

type conversationName struct {
	name string
}

func (n *conversationName) String() string {
	return n.name
}

func NewConversationName(name string) (*conversationName, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}

	if len(name) > 100 {
		return nil, errors.New("name is too long")
	}

	return &conversationName{
		name: name,
	}, nil
}

type GroupConversation struct {
	Conversation
	ID     uuid.UUID
	Name   conversationName
	Avatar string
	Owner  Participant
}

func NewGroupConversation(id uuid.UUID, name *conversationName, creator *User) (*GroupConversation, error) {
	groupConversation := &GroupConversation{
		Conversation: Conversation{
			ID:       id,
			Type:     ConversationTypeGroup,
			IsActive: true,
		},
		ID:     uuid.New(),
		Name:   *name,
		Avatar: string(name.String()[0]),
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
		return errors.New("user is not in conversation")
	}

	if groupConversation.Owner.UserID != participant.UserID {
		return errors.New("user is not owner")
	}

	if !groupConversation.IsActive {
		return errors.New("conversation is not active")
	}

	groupConversation.IsActive = false

	groupConversation.AddEvent(newGroupConversationDeletedEvent(groupConversation.Conversation.ID))

	return nil
}

func (groupConversation *GroupConversation) Rename(newName *conversationName, participant *Participant) error {
	if !groupConversation.isJoined(participant) {
		return errors.New("user is not in conversation")
	}

	if groupConversation.Owner.UserID != participant.UserID {
		return errors.New("user is not owner")
	}

	groupConversation.Name = *newName
	groupConversation.Avatar = string(newName.String()[0])

	groupConversation.AddEvent(newGroupConversationRenamedEvent(groupConversation.Conversation.ID, participant.UserID, newName.String()))

	return nil
}

func (groupConversation *GroupConversation) Join(user *User) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	participant := NewParticipant(uuid.New(), groupConversation.Conversation.ID, user.ID)

	participant.AddEvent(newGroupConversationJoinedEvent(groupConversation.Conversation.ID, user.ID))

	return participant, nil
}

func (groupConversation *GroupConversation) Leave(participant *Participant) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if !groupConversation.isJoined(participant) {
		return nil, errors.New("user is not in conversation")
	}

	participant.IsActive = false

	participant.AddEvent(newGroupConversationLeftEvent(groupConversation.Conversation.ID, participant.UserID))

	return participant, nil
}

func (groupConversation *GroupConversation) Invite(inviter *Participant, invitee *User) (*Participant, error) {
	if !groupConversation.isJoined(inviter) {
		return nil, errors.New("user is not in conversation")
	}

	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if groupConversation.Owner.UserID == invitee.ID {
		return nil, errors.New("user is owner")
	}

	if inviter.UserID == invitee.ID {
		return nil, errors.New("cannot invite yourself")
	}

	participant := NewParticipant(uuid.New(), groupConversation.Conversation.ID, invitee.ID)

	participant.AddEvent(newGroupConversationInvitedEvent(groupConversation.Conversation.ID, inviter.UserID, invitee.ID))

	return participant, nil
}

func (groupConversation *GroupConversation) Kick(kicker *Participant, target *Participant) (*Participant, error) {
	if !groupConversation.isJoined(target) {
		return nil, errors.New("user is not in conversation")
	}

	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if groupConversation.Owner.UserID != kicker.UserID {
		return nil, errors.New("user is not owner")
	}

	if kicker.UserID == target.UserID {
		return nil, errors.New("cannot kick yourself")
	}

	target.IsActive = false

	target.AddEvent(newGroupConversationLeftEvent(groupConversation.Conversation.ID, target.UserID))

	return target, nil
}

func (groupConversation *GroupConversation) SendTextMessage(messageID uuid.UUID, text string, participant *Participant) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if !groupConversation.isJoined(participant) {
		return nil, errors.New("user is not in conversation")
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
		return nil, errors.New("conversation is not active")
	}

	message := newJoinedConversationMessage(messageID, groupConversation.Conversation.ID, user.ID)

	return message, nil
}

func (groupConversation *GroupConversation) SendInvitedConversationMessage(messageID uuid.UUID, user *User) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	message := newInvitedConversationMessage(messageID, groupConversation.Conversation.ID, user.ID)

	return message, nil
}

func (groupConversation *GroupConversation) SendRenamedConversationMessage(messageID uuid.UUID, participant *Participant, newName string) (*Message, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	content := newRenamedMessageContent(newName)
	message := newConversationRenamedMessage(messageID, groupConversation.Conversation.ID, participant.UserID, content)

	return message, nil
}

func (groupConversation *GroupConversation) SendLeftConversationMessage(messageID uuid.UUID, participant *Participant) (*Message, error) {
	message := newLeftConversationMessage(messageID, groupConversation.Conversation.ID, participant.UserID)

	return message, nil
}
