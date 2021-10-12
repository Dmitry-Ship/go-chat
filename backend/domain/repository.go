package domain

import "github.com/google/uuid"

type RoomRepository interface {
	Create(room *Room) (*Room, error)
	Update(room *Room) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Room, error)
	FindByName(name string) (*Room, error)
	FindAll() ([]*Room, error)
}

type UserRepository interface {
	Create(user *User) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*User, error)
	FindByName(name string) (*User, error)
	FindAll() ([]*User, error)
}

type ChatMessageRepository interface {
	Create(message *ChatMessage) (*ChatMessage, error)
	Update(message *ChatMessage) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*ChatMessage, error)
	FindByRoomID(roomId uuid.UUID) ([]*ChatMessage, error)
	FindAll() ([]*ChatMessage, error)
}

type ParticipantRepository interface {
	Create(participant *Participant) (*Participant, error)
	Update(participant *Participant) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Participant, error)
	FindByRoomID(roomId uuid.UUID) ([]*Participant, error)
	FindByUserID(userId uuid.UUID) (*Participant, error)
	FindByRoomIDAndUserID(roomId uuid.UUID, userId uuid.UUID) (*Participant, error)
	FindAll() ([]*Participant, error)
	DeleteByRoomID(roomId uuid.UUID) error
}
