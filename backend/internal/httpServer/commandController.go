package httpServer

import (
	"GitHub/go-chat/backend/internal/app"
)

type CommandController struct {
	commands *app.Commands
}

func NewCommandController(commands *app.Commands) *CommandController {
	return &CommandController{
		commands: commands,
	}
}
