package domain

import (
	"fmt"
)

type subscription struct {
	user   *User
	roomId int64
}

// Room maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Rooms     map[int64]*Room
	Broadcast chan Message
	Join      chan subscription
	Leave     chan subscription
}

func NewHub() *Hub {
	return &Hub{
		Broadcast: make(chan Message, 1024),
		Join:      make(chan subscription),
		Leave:     make(chan subscription),
		Rooms:     make(map[int64]*Room),
	}
}

func NewSubscription(user *User, roomId int64) subscription {
	return subscription{
		user:   user,
		roomId: roomId,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case subscription := <-hub.Join:
			room := hub.Rooms[subscription.roomId]

			if room == nil {
				room = NewRoom("default")
				hub.Rooms[subscription.roomId] = room
			}
			user := subscription.user
			room.Users[user] = true

			message := NewMessage(fmt.Sprintf("%s %s joined", user.Avatar, user.Name), "system", room.Id, user)

			room.broadcastMessage(message)
		case subscription := <-hub.Leave:
			room := hub.Rooms[subscription.roomId]
			user := subscription.user

			if _, ok := room.Users[user]; ok {
				delete(room.Users, user)
				close(user.Send)
			}

			message := NewMessage(fmt.Sprintf("%s %s left", user.Avatar, user.Name), "system", room.Id, user)

			room.broadcastMessage(message)
		case message := <-hub.Broadcast:
			room := hub.Rooms[message.RoomId]
			room.broadcastMessage(message)
		}
	}
}
