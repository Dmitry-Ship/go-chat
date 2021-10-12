package application

type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type subscription struct {
	userId int32
	roomId int32
}

type Hub interface {
	NewNotification(notificationType string, data interface{}) Notification
	Register(client *Client)
	Unregister(client *Client)
	JoinRoom(userId int32, roomId int32)
	LeaveRoom(userId int32, roomId int32)
	DeleteRoom(roomId int32)
	Run()
}

type hub struct {
	Broadcast   chan *MessageFull
	deleteRoom  chan int32
	register    chan *Client
	unregister  chan *Client
	clients     map[int32][]*Client
	roomClients map[int32][]*Client
	joinRoom    chan subscription
	leaveRoom   chan subscription
}

func NewHub() *hub {
	return &hub{
		Broadcast:   make(chan *MessageFull, 1024),
		deleteRoom:  make(chan int32),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		joinRoom:    make(chan subscription),
		leaveRoom:   make(chan subscription),
		clients:     make(map[int32][]*Client),
		roomClients: make(map[int32][]*Client),
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
				RoomId int32 `json:"room_id"`
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

func (s *hub) JoinRoom(userId int32, roomId int32) {
	s.joinRoom <- subscription{
		userId: userId,
		roomId: roomId,
	}
}

func (s *hub) LeaveRoom(userId int32, roomId int32) {
	s.leaveRoom <- subscription{
		userId: userId,
		roomId: roomId,
	}
}

func (s *hub) DeleteRoom(roomId int32) {
	s.deleteRoom <- roomId
}

func (c *hub) NewNotification(notificationType string, data interface{}) Notification {
	return Notification{
		Type: notificationType,
		Data: data,
	}
}
