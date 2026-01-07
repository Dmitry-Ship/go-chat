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
	UserID  uuid.UUID   `json:"user_id"`
	Payload interface{} `json:"data"`
}

type IncomingNotification struct {
	Type   string
	Data   json.RawMessage
	UserID uuid.UUID
}

type connectionOptions struct {
	writeWait      time.Duration
	pongWait       time.Duration
	pingPeriod     time.Duration
	maxMessageSize int64
}

type Client struct {
	Id                         uuid.UUID
	UserID                     uuid.UUID
	connection                 *websocket.Conn
	sendChannel                chan OutgoingNotification
	handleIncomingNotification func(userID uuid.UUID, message []byte)
	unregisterClient           func(*Client)
	connectionOptions          connectionOptions
}

func NewClient(conn *websocket.Conn, unregisterClient func(*Client), handleIncomingNotification func(userID uuid.UUID, message []byte), userID uuid.UUID) *Client {
	return &Client{
		Id:                         uuid.New(),
		UserID:                     userID,
		connection:                 conn,
		sendChannel:                make(chan OutgoingNotification, 1024),
		unregisterClient:           unregisterClient,
		handleIncomingNotification: handleIncomingNotification,
		connectionOptions: connectionOptions{
			writeWait:      10 * time.Second,
			pongWait:       60 * time.Second,
			pingPeriod:     (60 * time.Second * 9) / 10,
			maxMessageSize: 512,
		},
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.unregisterClient(c)
		_ = c.connection.Close()
	}()

	c.connection.SetReadLimit(c.connectionOptions.maxMessageSize)
	err := c.connection.SetReadDeadline(time.Now().Add(c.connectionOptions.pongWait))

	if err != nil {
		log.Println(err)
	}

	c.connection.SetPongHandler(func(string) error {
		return c.connection.SetReadDeadline(time.Now().Add(c.connectionOptions.pongWait))
	})

	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleIncomingNotification(c.UserID, message)
	}
}

func (c *Client) WritePump() {
	pingTicker := time.NewTicker(c.connectionOptions.pingPeriod)
	defer func() {
		pingTicker.Stop()
		_ = c.connection.Close()
	}()

	for {
		select {
		case notification, ok := <-c.sendChannel:
			if err := c.connection.SetWriteDeadline(time.Now().Add(c.connectionOptions.writeWait)); err != nil {
				log.Println(err)
			}

			if !ok {
				// The conversation closed the channel.
				err := c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println(err)
				}
				return
			}

			if err := c.connection.WriteJSON(notification); err != nil {
				log.Println(err)
			}

		case <-pingTicker.C:
			if err := c.connection.SetWriteDeadline(time.Now().Add(c.connectionOptions.writeWait)); err != nil {
				log.Println(err)
			}

			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *Client) sendNotification(notification OutgoingNotification) {
	c.sendChannel <- notification
}
