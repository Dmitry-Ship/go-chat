package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
)

type HTTPHandlers struct {
	queryController   *QueryController
	commandController *CommandController
}

func NewHTTPHandlers(app *app.App) *HTTPHandlers {
	return &HTTPHandlers{
		queryController:   NewQueryController(app.Queries),
		commandController: NewCommandController(&app.Commands),
	}
}
