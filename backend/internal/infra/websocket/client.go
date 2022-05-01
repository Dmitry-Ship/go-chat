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
	Id                        uuid.UUID
	connection                *websocket.Conn
	unregisterChan            chan *client
	UserID                    uuid.UUID
	sendChannel               chan *OutgoingNotification
	incomingNotificationsChan chan *IncomingNotification
	connectionOptions         connectionOptions
}

func NewClient(conn *websocket.Conn, unregisterChan chan *client, incomingNotificationsChan chan *IncomingNotification, userID uuid.UUID) *client {
	return &client{
		Id:                        uuid.New(),
		UserID:                    userID,
		connection:                conn,
		sendChannel:               make(chan *OutgoingNotification, 1024),
		unregisterChan:            unregisterChan,
		incomingNotificationsChan: incomingNotificationsChan,
		connectionOptions: connectionOptions{
			writeWait:      10 * time.Second,
			pongWait:       60 * time.Second,
			pingPeriod:     (60 * time.Second * 9) / 10,
			maxMessageSize: 512,
		},
	}
}

func (c *client) ReadPump() {
	defer func() {
		c.unregisterChan <- c
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

		c.incomingNotificationsChan <- &IncomingNotification{
			Type:   notification.Type,
			Data:   data,
			UserID: c.UserID,
		}
	}
}

func (c *client) WritePump() {
	ticker := time.NewTicker(c.connectionOptions.pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
	}()

	for {
		select {
		case notification, ok := <-c.sendChannel:
			err := c.connection.SetWriteDeadline(time.Now().Add(c.connectionOptions.writeWait))

			if err != nil {
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

				return
			}

		case <-ticker.C:
			if err := c.connection.SetWriteDeadline(time.Now().Add(c.connectionOptions.writeWait)); err != nil {
				return
			}

			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *client) SendNotification(notification *OutgoingNotification) {
	c.sendChannel <- notification
}
