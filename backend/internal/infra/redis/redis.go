package pubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func GetRedisClient(ctx context.Context, conf RedisConfig) *redis.Client {
	options := fmt.Sprintf("%s:%s", conf.Host, conf.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     options,
		Password: conf.Password,
		DB:       0,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	log.Printf("📢 Connected to redis %s", options)

	return client
}

const ChatChannel = "chat"
const SubscriptionChannel = "subscriptions"
const PresenceChannel = "presence"
