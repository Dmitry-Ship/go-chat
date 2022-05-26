package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/readModel"
)

type HTTPHandlers struct {
	queryHandlers   *queryHandlers
	commandHandlers *commandHandlers
}

func NewHTTPHandlers(commands *app.Commands, queries readModel.QueriesRepository) *HTTPHandlers {
	return &HTTPHandlers{
		queryHandlers:   NewQueryHandlers(queries),
		commandHandlers: NewCommandHandlers(commands),
	}
}
