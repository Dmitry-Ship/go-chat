package server

import (
	"GitHub/go-chat/backend/internal/app"
	"GitHub/go-chat/backend/internal/domainEventsHandlers"
	"GitHub/go-chat/backend/internal/httpHandlers"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/readModel"
	"context"
)

type EventHandlers interface {
	ListenForEvents()
}

type Server struct {
	HttpHandlers  *httpHandlers.HTTPHandlers
	EventHandlers EventHandlers
}

func NewServer(ctx context.Context, commands *app.Commands, queries readModel.QueriesRepository, eventBus infra.EventsSubscriber) *Server {
	go commands.ClientsService.Run()

	return &Server{
		HttpHandlers:  httpHandlers.NewHTTPHandlers(commands, queries),
		EventHandlers: domainEventsHandlers.NewEventHandlers(ctx, eventBus, commands, queries),
	}
}
