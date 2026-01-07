package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheConfig struct {
	Prefix string
}

type CacheClient interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}

type redisCacheClient struct {
	client *redis.Client
	prefix string
}

func NewRedisCacheClient(redisClient *redis.Client, config CacheConfig) CacheClient {
	return &redisCacheClient{
		client: redisClient,
		prefix: config.Prefix,
	}
}

func (c *redisCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	fullKey := c.prefix + key
	data, err := c.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}
	return data, nil
}

func (c *redisCacheClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	fullKey := c.prefix + key
	err := c.client.Set(ctx, fullKey, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("cache set error: %w", err)
	}
	return nil
}

func (c *redisCacheClient) Delete(ctx context.Context, key string) error {
	fullKey := c.prefix + key
	err := c.client.Del(ctx, fullKey).Err()
	if err != nil {
		return fmt.Errorf("cache delete error: %w", err)
	}
	return nil
}

func (c *redisCacheClient) DeletePattern(ctx context.Context, pattern string) error {
	fullPattern := c.prefix + pattern
	iter := c.client.Scan(ctx, 0, fullPattern, 0).Iterator()
	keys := make([]string, 0)
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache scan error: %w", err)
	}
	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache delete pattern error: %w", err)
		}
	}
	return nil
}
