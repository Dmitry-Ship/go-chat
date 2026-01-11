package server

import (
	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/ratelimit"
	"GitHub/go-chat/backend/internal/readModel"
	"GitHub/go-chat/backend/internal/services"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (services.Tokens, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	SignUp(ctx context.Context, username string, password string) (services.Tokens, error)
	RotateTokens(ctx context.Context, refreshTokenString string) (services.Tokens, error)
	ParseAccessToken(accessTokenString string) (uuid.UUID, error)
}

type Server struct {
	ctx                  context.Context
	config               config.ServerConfig
	authCommands         AuthService
	groupConversation    services.GroupConversationService
	directConversation   services.DirectConversationService
	membership           services.MembershipService
	message              services.MessageService
	notificationCommands services.NotificationService
	queries              readModel.QueriesRepository
	ipRateLimiter        ratelimit.RateLimiter
	userRateLimiter      ratelimit.RateLimiter
}

func NewServer(
	ctx context.Context,
	config config.ServerConfig,
	authCommands AuthService,
	groupConversation services.GroupConversationService,
	directConversation services.DirectConversationService,
	membership services.MembershipService,
	message services.MessageService,
	notificationCommands services.NotificationService,
	queries readModel.QueriesRepository,
	ipRateLimiter ratelimit.RateLimiter,
	userRateLimiter ratelimit.RateLimiter,
) *Server {
	return &Server{
		ctx:                  ctx,
		config:               config,
		authCommands:         authCommands,
		groupConversation:    groupConversation,
		directConversation:   directConversation,
		membership:           membership,
		message:              message,
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
