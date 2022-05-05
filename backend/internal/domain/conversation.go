package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ConversationRepository interface {
	StorePublicConversation(conversation *PublicConversation) error
	StorePrivateConversation(conversation *PrivateConversation) error
	UpdatePublicConversation(conversation *PublicConversation) error
	GetPublicConversation(id uuid.UUID) (*PublicConversation, error)
	GetPrivateConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
}

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

func (publicConversation *PublicConversation) Delete(userID uuid.UUID) error {
	if publicConversation.Data.Owner.UserID != userID {
		return errors.New("user is not owner")
	}

	publicConversation.IsActive = false

	publicConversation.AddEvent(NewPublicConversationDeleted(publicConversation.Conversation.ID))

	return nil
}

func NewPublicConversation(id uuid.UUID, name string, creatorID uuid.UUID) *PublicConversation {
	publicConversation := &PublicConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "public",
			CreatedAt: time.Now(),
			IsActive:  true,
		},
		Data: PublicConversationData{
			ID:     uuid.New(),
			Name:   name,
			Avatar: string(name[0]),
			Owner:  *NewOwnerParticipant(id, creatorID),
		},
	}

	publicConversation.AddEvent(NewPublicConversationCreated(id, creatorID))

	return publicConversation
}

func (publicConversation *PublicConversation) Rename(newName string, userId uuid.UUID) error {
	if publicConversation.Data.Owner.UserID == userId {
		publicConversation.Data.Name = newName
		publicConversation.Data.Avatar = string(newName[0])

		publicConversation.AddEvent(NewPublicConversationRenamed(publicConversation.ID, userId, newName))
		return nil
	}

	return errors.New("user is not owner")
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
	privateConversation := PrivateConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "private",
			CreatedAt: time.Now(),
			IsActive:  true,
		},
		Data: PrivateConversationData{
			ID:       uuid.New(),
			ToUser:   *NewPrivateParticipant(id, to),
			FromUser: *NewPrivateParticipant(id, from),
		},
	}

	privateConversation.AddEvent(NewPrivateConversationCreated(id, to, from))

	return &privateConversation
}

func (privateConversation *PrivateConversation) GetFromUser() *Participant {
	return &privateConversation.Data.FromUser
}

func (privateConversation *PrivateConversation) GetToUser() *Participant {
	return &privateConversation.Data.ToUser
}
