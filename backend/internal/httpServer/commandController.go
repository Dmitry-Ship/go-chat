package httpServer

import (
	"GitHub/go-chat/backend/internal/app"
	ws "GitHub/go-chat/backend/internal/infra/websocket"
)

type CommandController struct {
	commands         *app.Commands
	wsConnectionsHub ws.Hub
}

func NewCommandController(commands *app.Commands, hub ws.Hub) *CommandController {
	return &CommandController{
		commands:         commands,
		wsConnectionsHub: hub,
	}
}
