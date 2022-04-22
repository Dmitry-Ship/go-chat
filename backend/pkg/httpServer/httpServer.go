package httpServer

import (
	"GitHub/go-chat/backend/pkg/app"
	ws "GitHub/go-chat/backend/pkg/websocket"
)

type HTTPServer struct {
	app              app.App
	WSConnectionsHub ws.Hub
}

func NewHTTPServer(app *app.App, hub ws.Hub) *HTTPServer {
	return &HTTPServer{
		app:              *app,
		WSConnectionsHub: hub,
	}
}
