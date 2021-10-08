package domain

import (
	"fmt"
)

type subscription struct {
	user   *User
	roomId int32
}

type Hub struct {
	Rooms     map[int32]*Room
	Broadcast chan ChatMessage
	Join      chan subscription
	Leave     chan subscription
}

func NewHub() *Hub {
	return &Hub{
		Broadcast: make(chan ChatMessage, 1024),
		Join:      make(chan subscription),
		Leave:     make(chan subscription),
		Rooms:     make(map[int32]*Room),
	}
}

func NewSubscription(user *User, roomId int32) subscription {
	return subscription{
		user:   user,
		roomId: roomId,
	}
}

func (hub *Hub) Run() {
	defaultRoom := NewRoom("default")
	hub.Rooms[defaultRoom.Id] = defaultRoom

	for {
		select {
		case subscription := <-hub.Join:
			room := hub.Rooms[subscription.roomId]

			user := subscription.user

			room.Users[user] = true

			message := NewChatMessage(fmt.Sprintf("%s %s joined", user.Avatar, user.Name), "system", room.Id, user)

			room.broadcastMessage(message)
		case subscription := <-hub.Leave:
			room := hub.Rooms[subscription.roomId]

			user := subscription.user

			delete(room.Users, user)

			message := NewChatMessage(fmt.Sprintf("%s %s left", user.Avatar, user.Name), "system", room.Id, user)

			room.broadcastMessage(message)
		case message := <-hub.Broadcast:
			room := hub.Rooms[message.RoomId]
			room.broadcastMessage(message)
		}
	}
}
