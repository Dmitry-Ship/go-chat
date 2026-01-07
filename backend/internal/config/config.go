package config

import "time"

type Token struct {
	Secret string
	TTL    time.Duration
}

type Auth struct {
	RefreshToken Token
	AccessToken  Token
}

type RateLimitConfig struct {
	MaxUserConnections int
	MaxIPConnections   int
	WindowDuration     time.Duration
}

type ServerConfig struct {
	Port         string
	ClientOrigin string
	RateLimit    RateLimitConfig
}
