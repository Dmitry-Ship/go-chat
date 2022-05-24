package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type BaseConversation interface {
	GetBaseData() *Conversation
}

type Conversation struct {
	aggregate
	ID        uuid.UUID
	Type      string
	CreatedAt time.Time
	IsActive  bool
}

func (c *Conversation) GetBaseData() *Conversation {
	return c
}

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

func (groupConversation *GroupConversation) Delete(userID uuid.UUID) error {
	if groupConversation.Data.Owner.UserID != userID {
		return errors.New("user is not owner")
	}

	groupConversation.IsActive = false

	groupConversation.AddEvent(NewGroupConversationDeleted(groupConversation.Conversation.ID))

	return nil
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

	groupConversation.AddEvent(NewGroupConversationCreated(id, creatorID))

	return groupConversation, nil
}

func (groupConversation *GroupConversation) Rename(newName string, userId uuid.UUID) error {
	if groupConversation.Data.Owner.UserID == userId {
		groupConversation.Data.Name = newName
		groupConversation.Data.Avatar = string(newName[0])

		groupConversation.AddEvent(NewGroupConversationRenamed(groupConversation.ID, userId, newName))
		return nil
	}

	return errors.New("user is not owner")
}

func (groupConversation *GroupConversation) SendTextMessage(text string, participant *Participant) (*TextMessage, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if participant.ConversationID != groupConversation.Conversation.ID {
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

	participant.AddEvent(NewGroupConversationJoined(groupConversation.ID, userID))

	return participant, nil
}

func (groupConversation *GroupConversation) Invite(inviteeID uuid.UUID) (*Participant, error) {
	if !groupConversation.Conversation.IsActive {
		return nil, errors.New("conversation is not active")
	}

	if groupConversation.Data.Owner.UserID == inviteeID {
		return nil, errors.New("user is owner")
	}

	participant := NewParticipant(groupConversation.ID, inviteeID)

	participant.AddEvent(NewGroupConversationInvited(groupConversation.ID, inviteeID))

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

type DirectConversationRepository interface {
	Store(conversation *DirectConversation) error
	GetID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*DirectConversation, error)
}

type DirectConversationData struct {
	ID       uuid.UUID
	ToUser   Participant
	FromUser Participant
}

type DirectConversation struct {
	Conversation
	Data DirectConversationData
}

func NewDirectConversation(id uuid.UUID, to uuid.UUID, from uuid.UUID) (*DirectConversation, error) {
	if to == from {
		return nil, errors.New("cannot chat with yourself")
	}

	directConversation := DirectConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "direct",
			CreatedAt: time.Now(),
			IsActive:  true,
		},
		Data: DirectConversationData{
			ID:       uuid.New(),
			ToUser:   *NewParticipant(id, to),
			FromUser: *NewParticipant(id, from),
		},
	}

	directConversation.AddEvent(NewDirectConversationCreated(id, to, from))

	return &directConversation, nil
}

func (directConversation *DirectConversation) GetFromUser() *Participant {
	return &directConversation.Data.FromUser
}

func (directConversation *DirectConversation) GetToUser() *Participant {
	return &directConversation.Data.ToUser
}

func (directConversation *DirectConversation) SendTextMessage(text string, userID uuid.UUID) (*TextMessage, error) {
	if directConversation.Data.ToUser.UserID != userID && directConversation.Data.FromUser.UserID != userID {
		return nil, errors.New("user is not participant")
	}

	message, err := newTextMessage(directConversation.ID, userID, text)

	if err != nil {
		return nil, err
	}

	return message, nil
}
