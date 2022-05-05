package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/readModel"
)

type HTTPHandlers struct {
	queryController   *queryController
	commandController *commandController
}

func NewHTTPHandlers(commands *app.Commands, queries readModel.QueriesRepository) *HTTPHandlers {
	return &HTTPHandlers{
		queryController:   NewQueryController(queries),
		commandController: NewCommandController(commands),
	}
}
