package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	GetURL(ctx context.Context, shortCode string) (string, error)
	SetURL(ctx context.Context, shortCode, longURL string) error
	Close() error
}

// RedisClient implements the Cache interface for Redis.
type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(redisUrl string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		client.Close()
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) GetURL(ctx context.Context, shortCode string) (string, error) {
	val, err := r.client.Get(ctx, "url:"+shortCode).Result()
	if err == redis.Nil {
		return "", nil // Cache miss is not an error
	}
	if err != nil {
		return "", err
	}
	return val, err
}

func (r *RedisClient) SetURL(ctx context.Context, shortCode, longURL string) error {
	err := r.client.Set(ctx, "url:"+shortCode, longURL, 24*time.Hour).Err()
	return err
}
