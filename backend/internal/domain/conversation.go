package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PrivateConversationRepository interface {
	Store(conversation *PrivateConversation) error
	GetID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*PrivateConversation, error)
}

type PublicConversationRepository interface {
	Store(conversation *PublicConversation) error
	Update(conversation *PublicConversation) error
	GetByID(id uuid.UUID) (*PublicConversation, error)
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

func NewPublicConversation(id uuid.UUID, name string, creatorID uuid.UUID) (*PublicConversation, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}

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
			Owner:  *NewParticipant(id, creatorID),
		},
	}

	publicConversation.AddEvent(NewPublicConversationCreated(id, creatorID))

	return publicConversation, nil
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

func (publicConversation *PublicConversation) SendTextMessage(text string, participant *Participant) (*TextMessage, error) {
	if participant.ConversationID != publicConversation.Conversation.ID {
		return nil, errors.New("user is not participant")
	}

	message, err := NewTextMessage(publicConversation.Conversation.ID, participant.UserID, text)

	if err != nil {
		return nil, err
	}

	return message, nil
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

func NewPrivateConversation(id uuid.UUID, to uuid.UUID, from uuid.UUID) (*PrivateConversation, error) {
	if to == from {
		return nil, errors.New("cannot chat with yourself")
	}

	privateConversation := PrivateConversation{
		Conversation: Conversation{
			ID:        id,
			Type:      "private",
			CreatedAt: time.Now(),
			IsActive:  true,
		},
		Data: PrivateConversationData{
			ID:       uuid.New(),
			ToUser:   *NewParticipant(id, to),
			FromUser: *NewParticipant(id, from),
		},
	}

	privateConversation.AddEvent(NewPrivateConversationCreated(id, to, from))

	return &privateConversation, nil
}

func (privateConversation *PrivateConversation) GetFromUser() *Participant {
	return &privateConversation.Data.FromUser
}

func (privateConversation *PrivateConversation) GetToUser() *Participant {
	return &privateConversation.Data.ToUser
}

func (privateConversation *PrivateConversation) SendTextMessage(text string, userID uuid.UUID) (*TextMessage, error) {
	if privateConversation.Data.ToUser.UserID != userID && privateConversation.Data.FromUser.UserID != userID {
		return nil, errors.New("user is not participant")
	}

	message, err := NewTextMessage(privateConversation.ID, userID, text)

	if err != nil {
		return nil, err
	}

	return message, nil
}
