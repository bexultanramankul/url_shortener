package usecase

import (
	"url_shortener/internal/repository"
	"url_shortener/internal/usecase/cache"
	"url_shortener/pkg/logger"
)

type UrlUsecase struct {
	urlRepository      repository.UrlRepository
	urlCacheRepository repository.UrlCacheRepository
	hashCache          cache.HashCache
}

func NewUrlUsecase(urlRepository repository.UrlRepository, urlCacheRepository repository.UrlCacheRepository, hashCache cache.HashCache) *UrlUsecase {
	return &UrlUsecase{urlRepository: urlRepository, urlCacheRepository: urlCacheRepository, hashCache: hashCache}
}

func (us *UrlUsecase) ShortenUrl(url string) (string, error) {
	hash, err := us.hashCache.GetHash()
	if err != nil {
		logger.Log.Errorf("Error getting hash from cache: %v", err)
		return "", err
	}

	err = us.urlRepository.Save(url, hash)
	if err != nil {
		logger.Log.Errorf("Error saving URL to repository: %v", err)
		return "", err
	}

	err = us.urlCacheRepository.Save(url, hash)
	if err != nil {
		logger.Log.Warnf("Error caching URL %s: %v", url, err)
	}

	return hash, nil
}

func (us *UrlUsecase) GetUrl(hash string) (string, error) {
	cachedUrl, err := us.urlCacheRepository.Get(hash)
	if err == nil && cachedUrl != "" {
		return cachedUrl, nil
	}

	url, err := us.urlRepository.FindUrlByHash(hash)
	if err != nil {
		logger.Log.Errorf("Error finding URL in repository for hash %s: %v", hash, err)
		return "", err
	}

	err = us.urlCacheRepository.Save(hash, url)
	if err != nil {
		logger.Log.Warnf("Error caching URL for hash %s: %v", hash, err)
	}

	return url, nil
}
