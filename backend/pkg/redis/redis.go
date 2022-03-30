package pubsub

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func GetRedisClient() *redis.Client {
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	options := fmt.Sprintf("%s:%s", host, port)
	fmt.Println(options)

	client := redis.NewClient(&redis.Options{
		Addr:     options,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Connected to redis %s", options))

	return client
}

const ChatChannel = "chat"
