package domain

import "github.com/google/uuid"

type RoomRepository interface {
	Store(room *Room) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Room, error)
	FindAll() ([]*Room, error)
}

type UserRepository interface {
	Store(user *User) error
	FindByID(id uuid.UUID) (*User, error)
}

type ChatMessageRepository interface {
	Store(message *ChatMessage) error
	FindAllByRoomID(roomId uuid.UUID) ([]*ChatMessage, error)
}

type ParticipantRepository interface {
	Store(participant *Participant) error
	DeleteByRoomIDAndUserID(roomId uuid.UUID, userId uuid.UUID) error
	FindAllByRoomID(roomId uuid.UUID) ([]*Participant, error)
	FindByRoomIDAndUserID(roomId uuid.UUID, userId uuid.UUID) (*Participant, error)
}
