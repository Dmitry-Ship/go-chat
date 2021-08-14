package application

import (
	"GitHub/go-chat/server/domain"
	"encoding/json"
	"fmt"
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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Id     string
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan domain.Message
	Sender *domain.Sender
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	id := strconv.Itoa(int(time.Now().UnixNano()))
	return &Client{
		Id:     id,
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan domain.Message, 256),
		Sender: domain.NewSender(id),
	}
}

// ReceiveMessages pumps messages from the websocket connection to the hub.
//
// The application runs ReceiveMessages in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReceiveMessages() {
	defer func() {
		c.Hub.Leave <- c
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

		m := domain.Message{}

		if err := json.Unmarshal(message, &m); err != nil {
			panic(err)
		}

		parsedMessage := domain.NewMessage(m.Content, "user", c.Sender)

		c.Hub.Broadcast <- parsedMessage
	}
}

// SendMessages pumps messages from the hub to the websocket connection.
//
// A goroutine running SendMessages is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) SendMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if message.Sender.Id == c.Id {
				message.Type = "outbound"
			}

			messages := []domain.Message{message}

			if err := c.Conn.WriteJSON(messages); err != nil {
				return
			}
			fmt.Println("sent")

			// Add queued chat messages to the current websocket message.
			for i := 0; i < len(c.Send); i++ {
				fmt.Println("in a loop")
				message := <-c.Send
				if message.Sender.Id == c.Id {
					message.Type = "outbound"
				}
				messages = append(messages, message)
				if err := c.Conn.WriteJSON(messages); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
