package server

import (
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	"context"
)

type EventHandlers interface {
	ListenForEvents()
}

type Server struct {
	ctx                  context.Context
	authCommands         services.AuthService
	conversationCommands services.ConversationService
	notificationCommands services.NotificationService
	queries              readModel.QueriesRepository
	subscriber           infra.EventsSubscriber
}

func NewServer(
	ctx context.Context,
	authCommands services.AuthService,
	conversationCommands services.ConversationService,
	notificationCommands services.NotificationService,
	queries readModel.QueriesRepository,
	eventBus infra.EventsSubscriber,
) *Server {
	return &Server{
		ctx:                  ctx,
		authCommands:         authCommands,
		conversationCommands: conversationCommands,
		notificationCommands: notificationCommands,
		queries:              queries,
		subscriber:           eventBus,
	}
}

func (s *Server) Run() {
	s.initRoutes()
	go s.listenForEvents()
	go s.notificationCommands.Run()
}
