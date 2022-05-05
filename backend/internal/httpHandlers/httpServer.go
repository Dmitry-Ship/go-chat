package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/hub"
	"GitHub/go-chat/backend/internal/readModel"
)

type HTTPHandlers struct {
	queryController   *queryController
	commandController *commandController
}

func NewHTTPHandlers(commands *app.Commands, queries readModel.QueriesRepository, clientRegister hub.ClientRegister) *HTTPHandlers {
	return &HTTPHandlers{
		queryController:   NewQueryController(queries),
		commandController: NewCommandController(commands, clientRegister),
	}
}
