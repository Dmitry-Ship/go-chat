package ws

import "time"

const (
	WriteWait       = 10 * time.Second
	PongWait        = 60 * time.Second
	PingPeriod      = (60 * time.Second * 9) / 10
	MaxMessageSize  = 512
	SendChannelSize = 1024
)
