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

type ConversationService interface {
	CreateGroupConversation(ctx context.Context, conversationID uuid.UUID, name string, userID uuid.UUID) error
	DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Rename(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, name string) error
	Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
	Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error
	StartDirectConversation(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error)
	SendTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error
}

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
	conversationCommands ConversationService
	notificationCommands services.NotificationService
	queries              readModel.QueriesRepository
	ipRateLimiter        ratelimit.RateLimiter
	userRateLimiter      ratelimit.RateLimiter
}

func NewServer(
	ctx context.Context,
	config config.ServerConfig,
	authCommands AuthService,
	conversationCommands ConversationService,
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
