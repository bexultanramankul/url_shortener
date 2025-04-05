package cache

import (
	"fmt"
	"sync"
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/repository"
	"url_shortener/internal/usecase/generator"
	"url_shortener/pkg/logger"
	"url_shortener/pkg/queue"
)

type HashCache interface {
	GetHash() (string, error)
}

type hashCache struct {
	hashRepository     repository.HashRepository
	uniqueIdRepository repository.UniqueIdRepository
	hashGenerator      generator.HashGenerator
	hashes             *queue.HashQueue
	mu                 sync.Mutex
	isRefilling        bool
	done               chan struct{}
}

func NewHashCache(
	hashRepository repository.HashRepository,
	uniqueIdRepository repository.UniqueIdRepository,
	hashGenerator generator.HashGenerator,
) HashCache {
	cache := &hashCache{
		hashRepository:     hashRepository,
		uniqueIdRepository: uniqueIdRepository,
		hashGenerator:      hashGenerator,
		hashes:             queue.NewHashQueue(config.AppConfig.Cache.Size),
		done:               make(chan struct{}),
	}
	cache.initFreeHashes()
	return cache
}

func (c *hashCache) initFreeHashes() {
	logger.Log.Info("Initializing hash cache...")

	err := c.checkAndGenerateHashesToDB()
	if err != nil {
		return
	}

	go c.fetchFreeHashes()
	logger.Log.Info("Hash cache initialized.")
}

func (c *hashCache) GetHash() (string, error) {
	for {
		select {
		case <-c.done:
			return "", fmt.Errorf("hash cache is shutting down")
		default:
			hash, ok := c.hashes.Pop()
			if ok {
				if c.hashes.Size() < config.AppConfig.Cache.RefillThreshold {
					logger.Log.Info("Cache size is below threshold, triggering refill...")
					go c.fetchFreeHashes()
				}
				return hash, nil
			}
			logger.Log.Info("Cache empty, waiting for refill...")
			time.Sleep(time.Duration(config.AppConfig.Cache.WaitTimeBeforeRetryMs) * time.Millisecond)
		}
	}
}

func (c *hashCache) fetchFreeHashes() {
	c.mu.Lock()
	if c.isRefilling {
		c.mu.Unlock()
		logger.Log.Debug("Refill already in progress, skipping fetch.")
		return
	}
	c.isRefilling = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.isRefilling = false
		c.mu.Unlock()
	}()

	logger.Log.Info("Fetching new hashes...")
	hashes, err := c.hashRepository.GetHashBatch(config.AppConfig.Cache.FetchBatchSize)
	if err != nil {
		logger.Log.Errorf("Failed to fetch new hashes from repository: %v", err)
		return
	}

	if len(hashes) == 0 {
		logger.Log.Warn("No new hashes fetched. Retrying in the next cycle...")
	} else {
		c.hashes.PushAll(hashes)
		logger.Log.Infof("Successfully added %d new hashes to the cache.", len(hashes))
	}

	go func() {
		err := c.checkAndGenerateHashesToDB()
		if err != nil {
			return
		}
	}()
}

func (c *hashCache) checkAndGenerateHashesToDB() error {
	logger.Log.Info("Checking hash count in the database...")
	count, err := c.hashRepository.GetHashCount()
	if err != nil {
		return fmt.Errorf("failed to get hash count from repository: %v", err)
	}

	if count < config.AppConfig.Cache.InitialHashCount {
		logger.Log.Infof("Hash count in database (%d) is less than required (%d), generating new hashes.", count, config.AppConfig.Cache.InitialHashCount)
		err := c.hashGenerator.GenerateHashBatch(config.AppConfig.Cache.InitialHashBatchSize)
		if err != nil {
			return fmt.Errorf("failed to generate hashes: %v", err)
		}
		logger.Log.Infof("Successfully generated %d new hashes.", config.AppConfig.Cache.InitialHashBatchSize)
	}

	return nil
}

func (c *hashCache) Shutdown() {
	close(c.done)
	logger.Log.Info("Hash cache shutting down...")
}
