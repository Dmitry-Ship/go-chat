package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type IncomingNotification struct {
	UserID uuid.UUID
	Data   json.RawMessage
}

type OutgoingNotification struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"data"`
}

type Client struct {
	Id             uuid.UUID
	hub            Hub
	wsHandler      WSHandlers
	Conn           *websocket.Conn
	send           chan *OutgoingNotification
	userID         uuid.UUID
	writeWait      time.Duration
	pongWait       time.Duration
	pingPeriod     time.Duration
	maxMessageSize int64
}

func NewClient(conn *websocket.Conn, hub Hub, wsHandler WSHandlers, userID uuid.UUID) *Client {
	// Time allowed to write a message to the peer.
	writeWait := 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait := 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod := (pongWait * 9) / 10
	var maxMessageSize int64 = 512

	return &Client{
		Id:             uuid.New(),
		userID:         userID,
		Conn:           conn,
		send:           make(chan *OutgoingNotification, 1024),
		hub:            hub,
		wsHandler:      wsHandler,
		writeWait:      writeWait,
		pongWait:       pongWait,
		pingPeriod:     pingPeriod,
		maxMessageSize: maxMessageSize,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(c.maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(c.pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(c.pongWait)); return nil })

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

		go c.wsHandler.HandleNotification(incomingNotification)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case notification, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if !ok {
				// The conversation closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(notification); err != nil {

				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendNotification(notification *OutgoingNotification) {
	c.send <- notification
}
