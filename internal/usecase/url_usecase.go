package usecase

import (
	"url_shortener/internal/repository"
	"url_shortener/internal/usecase/cache"
)

type UrlUsecase struct {
	urlRepository repository.UrlRepository
	hashCache     cache.HashCache
}

func NewUrlUsecase(urlRepository repository.UrlRepository, hashCache cache.HashCache) *UrlUsecase {
	return &UrlUsecase{urlRepository: urlRepository, hashCache: hashCache}
}

func (us *UrlUsecase) ShortenUrl(url string) (string, error) {
	hash, err := us.hashCache.GetHash()

	if err != nil {
		return "", err
	}

	err = us.urlRepository.Save(url, hash)
	if err != nil {
		return "", err
	}

	// TODO: Сохранение в Redis или другое кэширование

	return hash, nil
}

func (us *UrlUsecase) GetUrl(hash string) (string, error) {
	//TODO: сначала искать в редисе если нету то в бд

	url, err := us.urlRepository.FindUrlByHash(hash)
	if err != nil {
		return "", err
	}

	return url, nil
}
