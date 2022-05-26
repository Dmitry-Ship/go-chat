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
	httpHandlers               *httpHandlers.HTTPHandlers
	messageEventHandlers       EventHandlers
	notificationsEventHandlers EventHandlers
}

func NewServer(ctx context.Context, commands *app.Commands, queries readModel.QueriesRepository, eventBus infra.EventsSubscriber) *Server {
	go commands.ClientsService.Run()

	return &Server{
		httpHandlers:               httpHandlers.NewHTTPHandlers(commands, queries),
		messageEventHandlers:       domainEventsHandlers.NewMessageEventHandlers(ctx, eventBus, commands),
		notificationsEventHandlers: domainEventsHandlers.NewNotificationsEventHandlers(ctx, eventBus, commands, queries),
	}
}

func (s *Server) InitRoutes() {
	s.httpHandlers.InitRoutes()
}

func (s *Server) ListenForEvents() {
	go s.messageEventHandlers.ListenForEvents()
	go s.notificationsEventHandlers.ListenForEvents()
}
