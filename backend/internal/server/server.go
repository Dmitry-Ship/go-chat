package server

import (
	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/ratelimit"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	"context"
	"net/http"
)

type Server struct {
	ctx                  context.Context
	config               config.ServerConfig
	authCommands         services.AuthService
	conversationCommands services.ConversationService
	notificationCommands services.NotificationService
	queries              readModel.QueriesRepository
	ipRateLimiter        ratelimit.RateLimiter
	userRateLimiter      ratelimit.RateLimiter
}

func NewServer(
	ctx context.Context,
	config config.ServerConfig,
	authCommands services.AuthService,
	conversationCommands services.ConversationService,
	notificationCommands services.NotificationService,
	queries readModel.QueriesRepository,
	ipRateLimiter ratelimit.RateLimiter,
	userRateLimiter ratelimit.RateLimiter,
) *Server {
	return &Server{
		ctx:                  ctx,
		config:               config,
		authCommands:         authCommands,
		conversationCommands: conversationCommands,
		notificationCommands: notificationCommands,
		queries:              queries,
		ipRateLimiter:        ipRateLimiter,
		userRateLimiter:      userRateLimiter,
	}
}

func (s *Server) Run() http.Handler {
	mux := s.initRoutes()
	go s.notificationCommands.Run()
	return mux
}
