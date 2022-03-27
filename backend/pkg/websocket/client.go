package ws

import (
	"encoding/json"
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

type IncomingNotification struct {
	UserID uuid.UUID
	Data   json.RawMessage
}

type OutgoingNotification struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"data"`
	UserID  uuid.UUID   `json:"user_id"`
}

type Client struct {
	Id     uuid.UUID
	hub    Hub
	Conn   *websocket.Conn
	send   chan *OutgoingNotification
	userID uuid.UUID
}

func NewClient(conn *websocket.Conn, hub Hub, userID uuid.UUID) *Client {
	return &Client{
		Id:     uuid.New(),
		Conn:   conn,
		send:   make(chan *OutgoingNotification, 1024),
		userID: userID,
		hub:    hub,
	}
}

func (c *Client) ReceiveMessages() {
	defer func() {
		c.hub.UnregisterClient(c)
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

		incomingNotification := IncomingNotification{
			UserID: c.userID,
			Data:   message,
		}

		c.hub.SendMessage(incomingNotification)
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

func (c *Client) SendNotification(notification *OutgoingNotification) {
	c.send <- notification
}
