package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func ConnectRedis(ctx context.Context) (*Redis, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, fmt.Errorf("REDIS_URL is not set")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Redis{Client: client}, nil
}
