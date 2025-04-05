package repository

import (
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/storage"

	"github.com/redis/go-redis/v9"
)

type UrlCacheRepository interface {
	Get(hash string) (string, error)
	Save(hash, url string) error
	GetAndRefresh(hash string) (string, error)
}

type urlCacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewUrlCacheRepository() UrlCacheRepository {
	return &urlCacheRepository{
		client: storage.RedisClient,
		ttl:    time.Duration(config.AppConfig.Redis.TTLHours) * time.Hour,
	}
}

func (r *urlCacheRepository) Get(hash string) (string, error) {
	return r.client.Get(storage.RedisCtx, hash).Result()
}

func (r *urlCacheRepository) Save(hash, url string) error {
	return r.client.Set(storage.RedisCtx, hash, url, r.ttl).Err()
}

func (r *urlCacheRepository) GetAndRefresh(hash string) (string, error) {
	val, err := r.client.Get(storage.RedisCtx, hash).Result()
	if err != nil {
		return "", err
	}
	_ = r.client.Expire(storage.RedisCtx, hash, r.ttl).Err()
	return val, nil
}
