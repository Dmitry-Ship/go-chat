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
	ctx                         context.Context
	authCommands                services.AuthService
	conversationCommands        services.ConversationService
	notificationPipelineService services.NotificationsPipeline
	notificationCommands        services.NotificationService
	queries                     readModel.QueriesRepository
	subscriber                  infra.EventsSubscriber
}

func NewServer(
	ctx context.Context,
	authCommands services.AuthService,
	conversationCommands services.ConversationService,
	notificationPipelineService services.NotificationsPipeline,
	notificationCommands services.NotificationService,
	queries readModel.QueriesRepository,
	eventBus infra.EventsSubscriber,
) *Server {
	return &Server{
		ctx:                         ctx,
		authCommands:                authCommands,
		conversationCommands:        conversationCommands,
		notificationPipelineService: notificationPipelineService,
		notificationCommands:        notificationCommands,
		queries:                     queries,
		subscriber:                  eventBus,
	}
}

func (s *Server) Run() {
	s.initRoutes()
	s.listenForEvents()
	go s.notificationCommands.Run()
}
