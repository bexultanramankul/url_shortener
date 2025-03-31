package storage

import (
	"database/sql"
	"fmt"
	"sync"
	"url_shortener/internal/config"
	"url_shortener/pkg/logger"

	_ "github.com/lib/pq"
)

var (
	DB   *sql.DB
	once sync.Once
)

func InitDB() {
	log := logger.Log

	once.Do(func() {
		cfg := config.AppConfig.Database

		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
		)

		var err error
		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Database connection error: ", err)
		}

		if err = DB.Ping(); err != nil {
			log.Fatal("Database is unreachable: ", err)
		}

		log.Info("Connected to PostgreSQL")
	})
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			logger.Log.Warn("Error closing storage: ", err)
		} else {
			logger.Log.Info("Database connection closed")
		}
	}
}
