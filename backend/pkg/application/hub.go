package application

import (
	"github.com/google/uuid"
)

type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type subscription struct {
	userId uuid.UUID
	roomId uuid.UUID
}

type Hub interface {
	NewNotification(notificationType string, data interface{}) Notification
	Register(client *Client)
	Unregister(client *Client)
	JoinRoom(userId uuid.UUID, roomId uuid.UUID)
	LeaveRoom(userId uuid.UUID, roomId uuid.UUID)
	DeleteRoom(roomId uuid.UUID)
	Run()
}

type hub struct {
	Broadcast   chan *MessageFull
	deleteRoom  chan uuid.UUID
	register    chan *Client
	unregister  chan *Client
	clients     map[uuid.UUID][]*Client
	roomClients map[uuid.UUID][]*Client
	joinRoom    chan subscription
	leaveRoom   chan subscription
}

func NewHub() *hub {
	return &hub{
		Broadcast:   make(chan *MessageFull, 1024),
		deleteRoom:  make(chan uuid.UUID),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		joinRoom:    make(chan subscription),
		leaveRoom:   make(chan subscription),
		clients:     make(map[uuid.UUID][]*Client),
		roomClients: make(map[uuid.UUID][]*Client),
	}
}

func (s *hub) Run() {
	for {
		select {
		case message := <-s.Broadcast:
			notification := s.NewNotification("message", message)

			clients := s.roomClients[message.RoomId]

			for _, client := range clients {
				client.Send <- notification
			}

		case client := <-s.register:
			s.clients[client.userID] = append(s.clients[client.userID], client)

		case client := <-s.unregister:
			for i, c := range s.clients[client.userID] {
				if c == client {
					s.clients[client.userID] = append(s.clients[client.userID][:i], s.clients[client.userID][i+1:]...)
					break
				}
			}

		case sub := <-s.joinRoom:
			clients := s.clients[sub.userId]

			s.roomClients[sub.roomId] = append(s.roomClients[sub.roomId], clients...)

		case sub := <-s.leaveRoom:
			clients := s.roomClients[sub.roomId]

			for i, c := range clients {
				if c.userID == sub.userId {
					s.roomClients[sub.roomId] = append(clients[:i], clients[i+1:]...)
					break
				}
			}

		case roomId := <-s.deleteRoom:
			clients := s.roomClients[roomId]

			message := struct {
				RoomId uuid.UUID `json:"room_id"`
			}{
				RoomId: roomId,
			}

			for _, client := range clients {
				client.Send <- s.NewNotification("room_deleted", message)
			}

			delete(s.roomClients, roomId)
		}

	}
}

func (s *hub) Register(client *Client) {
	s.register <- client
}

func (s *hub) Unregister(client *Client) {
	s.unregister <- client
}

func (s *hub) JoinRoom(userId uuid.UUID, roomId uuid.UUID) {
	s.joinRoom <- subscription{
		userId: userId,
		roomId: roomId,
	}
}

func (s *hub) LeaveRoom(userId uuid.UUID, roomId uuid.UUID) {
	s.leaveRoom <- subscription{
		userId: userId,
		roomId: roomId,
	}
}

func (s *hub) DeleteRoom(roomId uuid.UUID) {
	s.deleteRoom <- roomId
}

func (c *hub) NewNotification(notificationType string, data interface{}) Notification {
	return Notification{
		Type: notificationType,
		Data: data,
	}
}
