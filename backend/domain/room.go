package domain

// Room maintains the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	Users map[*User]bool
	Name  string
	Id    int64
}

func NewRoom(name string) *Room {
	return &Room{
		// Id:    int64(time.Now().UnixNano()),
		Id: 123,

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
