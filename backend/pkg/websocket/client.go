package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type OutgoingNotification struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"data"`
}

type Client struct {
	Id             uuid.UUID
	connection     *websocket.Conn
	hub            Hub
	userID         uuid.UUID
	sendChannel    chan *OutgoingNotification
	wsHandler      WSHandlers
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
		connection:     conn,
		sendChannel:    make(chan *OutgoingNotification, 1024),
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
		c.connection.Close()
	}()

	c.connection.SetReadLimit(c.maxMessageSize)
	c.connection.SetReadDeadline(time.Now().Add(c.pongWait))
	c.connection.SetPongHandler(func(string) error { c.connection.SetReadDeadline(time.Now().Add(c.pongWait)); return nil })

	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var data json.RawMessage

		incomingNotification := struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}{
			Data: &data,
		}

		if err := json.Unmarshal(message, &incomingNotification); err != nil {
			log.Println(err)
			return
		}

		go c.wsHandler.HandleNotification(incomingNotification.Type, data, c.userID)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
	}()

	for {
		select {
		case notification, ok := <-c.sendChannel:
			c.connection.SetWriteDeadline(time.Now().Add(c.writeWait))
			if !ok {
				// The conversation closed the channel.
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.connection.WriteJSON(notification); err != nil {

				return
			}

		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendNotification(notification *OutgoingNotification) {
	c.sendChannel <- notification
}
