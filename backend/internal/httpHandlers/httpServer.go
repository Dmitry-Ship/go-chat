package httpHandlers

type HTTPHandlers struct {
	queryController   *QueryController
	commandController *CommandController
}

func NewHTTPHandlers(queryController *QueryController, commandController *CommandController) *HTTPHandlers {
	return &HTTPHandlers{
		queryController:   queryController,
		commandController: commandController,
	}
}
