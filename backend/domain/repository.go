package domain

type RoomRepository interface {
	Create(room *Room) (*Room, error)
	Update(room *Room) error
	Delete(id int32) error
	FindByID(id int32) (*Room, error)
	FindByName(name string) (*Room, error)
	FindAll() ([]*Room, error)
}

type UserRepository interface {
	Create(user *User) (*User, error)
	Update(user *User) error
	Delete(id int32) error
	FindByID(id int32) (*User, error)
	FindByName(name string) (*User, error)
	FindAll() ([]*User, error)
}

type ChatMessageRepository interface {
	Create(message *ChatMessage) (*ChatMessage, error)
	Update(message *ChatMessage) error
	Delete(id int32) error
	FindByID(id int32) (*ChatMessage, error)
	FindByRoomID(roomID int32) ([]*ChatMessage, error)
	FindAll() ([]*ChatMessage, error)
}

type ParticipantRepository interface {
	Create(participant *Participant) (*Participant, error)
	Update(participant *Participant) error
	Delete(id int32) error
	FindByID(id int32) (*Participant, error)
	FindByRoomID(roomID int32) ([]*Participant, error)
	FindByUserID(userID int32) (*Participant, error)
	FindByRoomIDAndUserID(roomID int32, userID int32) (*Participant, error)
	FindAll() ([]*Participant, error)
}
