package redis

import (
	"context"
	"os"

	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func OpenRedisConnection(ctx context.Context, cfg *config.Config) *RedisClient {
	log := logger.GetLogger()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: cfg.Redis.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Unable to ping Redis connection", "error", err.Error())
		os.Exit(1)
	}

	return &RedisClient{
		Client: rdb,
	}
}
