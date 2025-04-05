package repository

import (
	"database/sql"
)

type UrlRepository interface {
	Save(url string, hash string) error
	FindUrlByHash(hash string) (string, error)
}

type urlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) UrlRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Save(url string, hash string) error {
	const query = `INSERT INTO url (hash, url) VALUES ($1, $2);`

	_, err := r.db.Exec(query, hash, url)
	if err != nil {
		return err
	}

	return nil
}

func (r *urlRepository) FindUrlByHash(hash string) (string, error) {
	const query = `SELECT url FROM url WHERE hash = $1;`
	var url string
	err := r.db.QueryRow(query, hash).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return url, nil
}
