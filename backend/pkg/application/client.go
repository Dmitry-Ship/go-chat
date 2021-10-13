package application

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type incomingNotification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type outgoingNotification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	Id             uuid.UUID
	messageService MessageService
	roomService    RoomService
	userService    UserService
	Conn           *websocket.Conn
	send           chan outgoingNotification
	userID         uuid.UUID
}

func NewClient(conn *websocket.Conn, messageService MessageService, userService UserService, roomService RoomService, userID uuid.UUID) *Client {
	return &Client{
		Id:             uuid.New(),
		messageService: messageService,
		roomService:    roomService,
		userService:    userService,
		Conn:           conn,
		send:           make(chan outgoingNotification, 1024),
		userID:         userID,
	}
}

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

		notification := incomingNotification{
			Data: &data,
		}

		if err := json.Unmarshal(message, &notification); err != nil {
			fmt.Println(err)
			continue
		}

		switch notification.Type {
		case "message":
			request := struct {
				Content string    `json:"content"`
				RoomId  uuid.UUID `json:"room_id"`
				UserId  uuid.UUID `json:"user_id"`
			}{}

			if err := json.Unmarshal([]byte(data), &request); err != nil {
				fmt.Println(err)
				return
			}

			_, err := c.messageService.SendMessage(request.Content, "user", request.RoomId, request.UserId)
			if err != nil {
				fmt.Println(err)
			}
		case "join":
			request := struct {
				RoomId uuid.UUID `json:"room_id"`
				UserId uuid.UUID `json:"user_id"`
			}{}

			if err := json.Unmarshal(data, &request); err != nil {
				fmt.Println(err)
			}

			_, err := c.roomService.JoinRoom(request.UserId, request.RoomId)
			if err != nil {
				fmt.Println(err)
			}
		case "leave":
			request := struct {
				RoomId uuid.UUID `json:"room_id"`
				UserId uuid.UUID `json:"user_id"`
			}{}

			if err := json.Unmarshal(data, &request); err != nil {
				fmt.Println(err)
				return
			}

			err := c.roomService.LeaveRoom(request.UserId, request.RoomId)

			if err != nil {
				fmt.Println(err)
			}

		}

	}
}

func (c *Client) SendNotifications() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case notification, ok := <-c.send:
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

func (c *Client) SendNotification(notificationType string, data interface{}) {
	c.send <- outgoingNotification{
		Type: notificationType,
		Data: data,
	}
}
