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
	ID        uuid.UUID
	Type      string
	CreatedAt time.Time
}

func (c *Conversation) GetBaseData() *Conversation {
	return c
}

type PublicConversationData struct {
	ID     uuid.UUID
	Name   string
	Avatar string
	Owner  Participant
}
type PublicConversation struct {
	Conversation
	Data PublicConversationData
}

func NewPublicConversation(id uuid.UUID, name string, creatorID uuid.UUID) *PublicConversation {
	return &PublicConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "public",
			CreatedAt: time.Now(),
		},
		Data: PublicConversationData{
			ID:     uuid.New(),
			Name:   name,
			Avatar: string(name[0]),
			Owner:  *NewParticipant(id, creatorID),
		},
	}
}

func (publicConversation *PublicConversation) Rename(newName string, userId uuid.UUID) error {
	if publicConversation.Data.Owner.UserID == userId {
		publicConversation.Data.Name = newName
		publicConversation.Data.Avatar = string(newName[0])
	} else {
		return errors.New("user is not owner")
	}
	return nil
}

type PrivateConversationData struct {
	ID       uuid.UUID
	ToUser   Participant
	FromUser Participant
}

type PrivateConversation struct {
	Conversation
	Data PrivateConversationData
}

func NewPrivateConversation(id uuid.UUID, to uuid.UUID, from uuid.UUID) *PrivateConversation {
	return &PrivateConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "private",
			CreatedAt: time.Now(),
		},
		Data: PrivateConversationData{
			ID:       uuid.New(),
			ToUser:   *NewParticipant(id, to),
			FromUser: *NewParticipant(id, from),
		},
	}
}

func (privateConversation *PrivateConversation) GetFromUser() *Participant {
	return &privateConversation.Data.FromUser
}

func (privateConversation *PrivateConversation) GetToUser() *Participant {
	return &privateConversation.Data.ToUser
}
