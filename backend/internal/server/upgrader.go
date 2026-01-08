package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	WebSocketBufferSize = 1024
)

var WebSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  WebSocketBufferSize,
	WriteBufferSize: WebSocketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
