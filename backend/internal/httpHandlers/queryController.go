package httpHandlers

import (
	"GitHub/go-chat/backend/internal/readModel"
)

type QueryController struct {
	queries readModel.QueriesRepository
}

func NewQueryController(queries readModel.QueriesRepository) *QueryController {
	return &QueryController{
		queries: queries,
	}
}
