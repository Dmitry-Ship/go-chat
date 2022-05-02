package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
)

type HTTPHandlers struct {
	queryController   *queryController
	commandController *commandController
}

func NewHTTPHandlers(app *app.App) *HTTPHandlers {
	return &HTTPHandlers{
		queryController:   NewQueryController(app.Queries),
		commandController: NewCommandController(&app.Commands),
	}
}
