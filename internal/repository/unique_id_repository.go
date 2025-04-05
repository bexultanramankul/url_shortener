package repository

import (
	"database/sql"
)

type UniqueIdRepository interface {
	GetUniqueNumbers(count int) ([]int64, error)
}

type uniqueIdRepository struct {
	db *sql.DB
}

func NewUniqueIdRepository(db *sql.DB) UniqueIdRepository {
	return &uniqueIdRepository{db: db}
}

func (r *uniqueIdRepository) GetUniqueNumbers(count int) ([]int64, error) {
	query := `SELECT nextval('unique_number_seq') FROM generate_series(1, $1);`
	rows, err := r.db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	numbers := make([]int64, count)
	i := 0
	for rows.Next() {
		if err := rows.Scan(&numbers[i]); err != nil {
			return nil, err
		}
		i++
	}
	numbers = numbers[:i]

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return numbers, nil
}
