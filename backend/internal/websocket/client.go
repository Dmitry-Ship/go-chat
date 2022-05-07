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
type client struct {
	Id                         uuid.UUID
	connection                 *websocket.Conn
	UserID                     uuid.UUID
	sendChannel                chan *OutgoingNotification
	handleincomingNotification func(notification *IncomingNotification)
	unregisterClient           func(client *client)
	connectionOptions          connectionOptions
}

func NewClient(conn *websocket.Conn, unregisterClient func(client *client), handleincomingNotification func(notification *IncomingNotification), userID uuid.UUID) *client {
	return &client{
		Id:                         uuid.New(),
		UserID:                     userID,
		connection:                 conn,
		sendChannel:                make(chan *OutgoingNotification, 1024),
		unregisterClient:           unregisterClient,
		handleincomingNotification: handleincomingNotification,
		connectionOptions: connectionOptions{
			writeWait:      10 * time.Second,
			pongWait:       60 * time.Second,
			pingPeriod:     (60 * time.Second * 9) / 10,
			maxMessageSize: 512,
		},
	}
}

func (c *client) Listen() {
	go c.writePump()
	go c.readPump()
}

func (c *client) readPump() {
	defer func() {
		go c.unregisterClient(c)
		c.connection.Close()
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

		var data json.RawMessage

		notification := struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}{
			Data: &data,
		}

		if err := json.Unmarshal(message, &notification); err != nil {
			log.Println(err)
			return
		}

		go c.handleincomingNotification(&IncomingNotification{
			Type:   notification.Type,
			Data:   data,
			UserID: c.UserID,
		})
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(c.connectionOptions.pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
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

		case <-ticker.C:
			if err := c.connection.SetWriteDeadline(time.Now().Add(c.connectionOptions.writeWait)); err != nil {
				log.Println(err)
			}

			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *client) sendNotification(notification *OutgoingNotification) {
	c.sendChannel <- notification
}
