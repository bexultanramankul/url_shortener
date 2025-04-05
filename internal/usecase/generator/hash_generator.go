package generator

import (
	"url_shortener/internal/pkg/encoder"
	"url_shortener/internal/repository"
)

type HashGenerator interface {
	GenerateHashBatch(count int) error
}

type hashGenerator struct {
	uniqueIDRepository repository.UniqueIdRepository
	hashRepository     repository.HashRepository
}

func NewHashGenerator(
	uniqueIDRepository repository.UniqueIdRepository,
	hashRepository repository.HashRepository,
) HashGenerator {
	return &hashGenerator{
		uniqueIDRepository: uniqueIDRepository,
		hashRepository:     hashRepository,
	}
}

func (hg *hashGenerator) GenerateHashBatch(count int) error {
	numbers, err := hg.uniqueIDRepository.GetUniqueNumbers(count)
	if err != nil {
		return err
	}

	hashes := make([]string, 0, count)
	for _, num := range numbers {
		hashes = append(hashes, encoder.Encode(num))
	}

	err = hg.hashRepository.SaveHashBatch(hashes)
	if err != nil {
		return err
	}

	return nil
}
