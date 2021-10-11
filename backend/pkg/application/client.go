package application

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type IncomingNotification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client is a middleman between the websocket connection and the room.
type Client struct {
	Id             string
	messageService MessageService
	roomService    RoomService
	userService    UserService
	Conn           *websocket.Conn
	Send           chan Notification
}

func NewClient(conn *websocket.Conn, messageService MessageService, userService UserService, roomService RoomService, Send chan Notification) *Client {
	id := strconv.Itoa(int(time.Now().UnixNano()))
	return &Client{
		Id:             id,
		messageService: messageService,
		roomService:    roomService,
		userService:    userService,
		Conn:           conn,
		Send:           Send,
	}
}

// ReceiveMessages pumps messages from the websocket connection to the room.
//
// The application runs ReceiveMessages in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func (c *Client) ReceiveMessages() {
	defer func() {
		// c.Hub.Leave <- domain.NewSubscription(c.User, 123)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var data json.RawMessage

		notification := IncomingNotification{
			Data: &data,
		}

		if err := json.Unmarshal(message, &notification); err != nil {
			panic(err)
		}

		switch notification.Type {
		case "message":
			request := struct {
				Content string `json:"content"`
				RoomId  int32  `json:"room_id"`
				UserId  int32  `json:"user_id"`
			}{}

			if err := json.Unmarshal([]byte(data), &request); err != nil {
				panic(err)
			}

			c.messageService.SendMessage(request.Content, "user", request.RoomId, request.UserId)
		case "join":
			request := struct {
				RoomId int32 `json:"room_id"`
				UserId int32 `json:"user_id"`
			}{}

			if err := json.Unmarshal(data, &request); err != nil {
				panic(err)
			}

			c.roomService.JoinRoom(request.UserId, request.RoomId)
		case "leave":
			request := struct {
				RoomId int32 `json:"room_id"`
				UserId int32 `json:"user_id"`
			}{}

			if err := json.Unmarshal(data, &request); err != nil {
				panic(err)
			}

			c.roomService.LeaveRoom(request.UserId, request.RoomId)
		}

	}
}

// SendNotifications pumps messages from the hub to the websocket connection.
//
// A goroutine running SendNotifications is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) SendNotifications() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case notification, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(notification); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
