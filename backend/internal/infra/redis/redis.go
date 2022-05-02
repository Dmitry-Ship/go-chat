package pubsub

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func GetRedisClient(ctx context.Context) *redis.Client {
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	options := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     options,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("📢 Connected to redis %s", options))

	return client
}

const ChatChannel = "chat"
