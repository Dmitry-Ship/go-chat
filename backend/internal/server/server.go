package server

import (
	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/infra"
	"GitHub/go-chat/backend/internal/ratelimit"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	"context"
)

type EventHandlers interface {
	ListenForEvents()
}

type Server struct {
	ctx                         context.Context
	config                      config.ServerConfig
	authCommands                services.AuthService
	conversationCommands        services.ConversationService
	notificationPipelineService services.NotificationsPipeline
	notificationCommands        services.NotificationService
	queries                     readModel.QueriesRepository
	subscriber                  infra.EventsSubscriber
	ipRateLimiter               ratelimit.RateLimiter
	userRateLimiter             ratelimit.RateLimiter
}

func NewServer(
	ctx context.Context,
	config config.ServerConfig,
	authCommands services.AuthService,
	conversationCommands services.ConversationService,
	notificationPipelineService services.NotificationsPipeline,
	notificationCommands services.NotificationService,
	queries readModel.QueriesRepository,
	eventBus infra.EventsSubscriber,
	ipRateLimiter ratelimit.RateLimiter,
	userRateLimiter ratelimit.RateLimiter,
) *Server {
	return &Server{
		ctx:                         ctx,
		config:                      config,
		authCommands:                authCommands,
		conversationCommands:        conversationCommands,
		notificationPipelineService: notificationPipelineService,
		notificationCommands:        notificationCommands,
		queries:                     queries,
		subscriber:                  eventBus,
		ipRateLimiter:               ipRateLimiter,
		userRateLimiter:             userRateLimiter,
	}
}

func (s *Server) Run() {
	s.initRoutes()
	s.listenForEvents()
	go s.notificationCommands.Run()
}
