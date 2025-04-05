package repository

import (
	"database/sql"
	"fmt"
	"url_shortener/pkg/logger"
)

type HashRepository interface {
	GetHash() (string, error)
	GetHashBatch(count int) ([]string, error)
	SaveHashBatch(hashes []string) error
	GetHashCount() (int, error)
}

type hashRepository struct {
	db *sql.DB
}

func NewHashRepository(db *sql.DB) HashRepository {
	return &hashRepository{db: db}
}

func (r *hashRepository) GetHash() (string, error) {
	const query = `
		DELETE FROM hash
		WHERE hash = (SELECT hash FROM hash LIMIT 1)
		RETURNING hash;
	`

	var hash string
	err := r.db.QueryRow(query).Scan(&hash)
	if err != nil {
		logger.Log.Errorf("Failed to get hash from database: %v", err)
		return "", fmt.Errorf("internal server error")
	}

	return hash, nil
}

func (r *hashRepository) GetHashBatch(count int) ([]string, error) {
	const query = `
		DELETE FROM hash
		WHERE hash IN (SELECT hash FROM hash LIMIT $1)
		RETURNING hash;
	`

	rows, err := r.db.Query(query, count)
	if err != nil {
		logger.Log.Errorf("Failed to query hashes batch: %v", err)
		return nil, fmt.Errorf("internal server error")
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.Log.Errorf("Error closing rows: %v", closeErr)
		}
	}()

	hashes := make([]string, count)
	i := 0
	for rows.Next() {
		if err := rows.Scan(&hashes[i]); err != nil {
			logger.Log.Errorf("Failed to scan hash: %v", err)
			return nil, fmt.Errorf("internal server error")
		}
		i++
	}
	hashes = hashes[:i]

	if err := rows.Err(); err != nil {
		logger.Log.Errorf("Row iteration error: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return hashes, nil
}

func (r *hashRepository) SaveHashBatch(hashes []string) error {
	if len(hashes) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Errorf("Failed to begin transaction: %v", err)
		return fmt.Errorf("internal server error")
	}

	stmt, err := tx.Prepare("INSERT INTO hash (hash) VALUES ($1)")
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logger.Log.Errorf("Failed to rollback transaction: %v", rollbackErr)
		}
		logger.Log.Errorf("Failed to prepare statement: %v", err)
		return fmt.Errorf("internal server error")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			logger.Log.Errorf("Error closing statement: %v", closeErr)
		}
	}()

	for _, hash := range hashes {
		_, err := stmt.Exec(hash)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logger.Log.Errorf("Failed to rollback transaction: %v", rollbackErr)
			}
			logger.Log.Errorf("Failed to execute insert for hash %s: %v", hash, err)
			return fmt.Errorf("internal server error")
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("internal server error")
	}

	return nil
}

func (r *hashRepository) GetHashCount() (int, error) {
	const query = `SELECT COUNT(*) FROM hash;`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		logger.Log.Errorf("Failed to get hash count from database: %v", err)
		return 0, fmt.Errorf("internal server error")
	}

	return count, nil
}
