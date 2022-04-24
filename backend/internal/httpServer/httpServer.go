package httpServer

type HTTPServer struct {
	queryController   *QueryController
	commandController *CommandController
}

func NewHTTPServer(queryController *QueryController, commandController *CommandController) *HTTPServer {
	return &HTTPServer{
		queryController:   queryController,
		commandController: commandController,
	}
}
