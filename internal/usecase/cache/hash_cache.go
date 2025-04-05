package cache

import (
	"fmt"
	"sync"
	"time"
	"url_shortener/internal/repository"
	"url_shortener/internal/usecase/generator"
	"url_shortener/pkg/logger"
	"url_shortener/pkg/queue"
)

const cacheSize = 10
const refillThreshold = 2
const fetchBatchSize = 10

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
		hashes:             queue.NewHashQueue(cacheSize),
		done:               make(chan struct{}),
	}
	cache.initFreeHashes()
	return cache
}

func (c *hashCache) initFreeHashes() {
	logger.Log.Info("Initializing hash cache...")

	// Генерация хешей
	err := c.hashGenerator.GenerateHashBatch(cacheSize * 2)
	if err != nil {
		logger.Log.Info("Failed to generate hashes: %v", err)
		return
	}

	// Запуск потока для пополнения хешей
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
				if c.hashes.Size() < refillThreshold {
					go c.fetchFreeHashes()
				}
				return hash, nil
			}
			logger.Log.Info("Cache empty, waiting for refill...")
			time.Sleep(100 * time.Millisecond)

			//TODO: repo getHash
		}
	}
}

func (c *hashCache) fetchFreeHashes() {
	c.mu.Lock()
	if c.isRefilling {
		c.mu.Unlock()
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
	hashes, err := c.hashRepository.GetHashBatch(fetchBatchSize)
	if err != nil {
		logger.Log.Info("Failed to fetch new hashes: %v", err)
		return
	}

	c.hashes.PushAll(hashes)
	logger.Log.Info("Added %d new hashes to cache.", len(hashes))

	//TODO: generator hashes if there is not enough hash in the database
}

func (c *hashCache) Shutdown() {
	close(c.done)
	logger.Log.Info("Hash cache shutting down...")
}
