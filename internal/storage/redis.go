package storage

import (
	"context"
	"url_shortener/internal/config"
	"url_shortener/pkg/logger"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	RedisCtx    = context.Background()
)

func InitRedis() {
	cfg := config.AppConfig.Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		logger.Log.Fatalf("Redis connection error: %v", err)
	}

	logger.Log.Info("Redis connected successfully")
}

func CloseRedis() {
	if err := RedisClient.Close(); err != nil {
		logger.Log.Warn("Error closing Redis: ", err)
	} else {
		logger.Log.Info("Redis connection closed")
	}
}
