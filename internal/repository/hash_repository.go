package repository

import (
	"database/sql"
)

type HashRepository interface {
	GetHash() (string, error)
	GetHashBatch(count int) ([]string, error)
	SaveHashBatch(hashes []string) error
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
		return "", err
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
		return nil, err
	}
	defer rows.Close()

	hashes := make([]string, count)
	i := 0
	for rows.Next() {
		if err := rows.Scan(&hashes[i]); err != nil {
			return nil, err
		}
		i++
	}
	hashes = hashes[:i]

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hashes, nil
}

func (r *hashRepository) SaveHashBatch(hashes []string) error {
	if len(hashes) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO hash (hash) VALUES ($1)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, hash := range hashes {
		_, err := stmt.Exec(hash)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
