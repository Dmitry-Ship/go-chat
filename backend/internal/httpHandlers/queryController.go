package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
)

type QueryController struct {
	queries *app.Queries
}

func NewQueryController(queries *app.Queries) *QueryController {
	return &QueryController{
		queries: queries,
	}
}
