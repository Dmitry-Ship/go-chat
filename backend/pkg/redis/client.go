package redis

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis"
)

type RedisClient interface {
	ReadFromChannel(channel string)
	SendToChannel(channel string, msg interface{})
	GetMessageChannel() <-chan string
}

type RedisConnection struct {
	client          *redis.Client
	IncomingMessage chan string
}

func NewRedisConnection() *RedisConnection {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("missing REDIS_HOST env var")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Fatal("missing REDIS_PASSWORD env var")
	}

	log.Println("connecting to Redis...")
	client := redis.NewClient(&redis.Options{Addr: redisHost, Password: redisPassword})

	_, err := client.Ping().Result()

	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}

	return &RedisConnection{
		client:          client,
		IncomingMessage: make(chan string, 1024),
	}
}

func (c *RedisConnection) ReadFromChannel(channel string) {
	sub := c.client.Subscribe(channel)
	messages := sub.Channel()

	for message := range messages {
		c.IncomingMessage <- message.Payload
	}
}

func (c *RedisConnection) GetMessageChannel() <-chan string {
	return c.IncomingMessage
}

func (c *RedisConnection) SendToChannel(channel string, msg interface{}) {
	payload, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	err = c.client.Publish(channel, payload).Err()

	if err != nil {
		log.Println("could not publish to channel", err)
	}
}
