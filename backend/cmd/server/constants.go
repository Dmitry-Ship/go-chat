package main

import "time"

const (
	DefaultAccessTokenTTL  = 10 * time.Minute
	DefaultRefreshTokenTTL = 24 * 90 * time.Hour
	DefaultUserRateLimit   = 10
	DefaultIPRateLimit     = 20
	DefaultRateLimitWindow = 60 * time.Second
)
