package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/readModel"
	ws "GitHub/go-chat/backend/internal/websocket"
)

type HTTPHandlers struct {
	queryController   *queryController
	commandController *commandController
}

func NewHTTPHandlers(commands *app.Commands, queries readModel.QueriesRepository, clientRegister ws.ClientRegister) *HTTPHandlers {
	return &HTTPHandlers{
		queryController:   NewQueryController(queries),
		commandController: NewCommandController(commands, clientRegister),
	}
}
