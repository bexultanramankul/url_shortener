package generator

import (
	"url_shortener/internal/pkg/encoder"
	"url_shortener/internal/repository"
	"url_shortener/pkg/logger"
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
	logger.Log.Infof("Generating batch of %d hashes...", count)

	numbers, err := hg.uniqueIDRepository.GetUniqueNumbers(count)
	if err != nil {
		logger.Log.Errorf("Failed to get unique numbers: %v", err)
		return err
	}

	logger.Log.Infof("Successfully retrieved %d unique numbers for hash generation.", len(numbers))

	hashes := make([]string, 0, count)
	for _, num := range numbers {
		hash := encoder.Encode(num)
		hashes = append(hashes, hash)
	}

	logger.Log.Infof("Successfully generated %d hashes.", len(hashes))

	err = hg.hashRepository.SaveHashBatch(hashes)
	if err != nil {
		logger.Log.Errorf("Failed to save hash batch: %v", err)
		return err
	}

	logger.Log.Infof("Successfully saved %d hashes to the repository.", len(hashes))

	return nil
}
