package domain

import "math/rand"

type Room struct {
	Users map[*User]bool `json:"-"`
	Name  string         `json:"name"`
	Id    int32          `json:"id"`
}

func NewRoom(name string) *Room {
	return &Room{
		Id: int32(rand.Int31()),

		Users: make(map[*User]bool),

		Name: name,
	}
}

func (room *Room) broadcastMessage(message interface{}) {
	for user := range room.Users {
		notification := user.NewNotification("message", message)

		select {
		case user.Send <- notification:
		default:
			close(user.Send)
			delete(room.Users, user)
		}
	}
}
