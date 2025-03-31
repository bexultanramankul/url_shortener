package main

import (
	"fmt"
	"net/http"
	"url_shortener/internal/config"
	"url_shortener/internal/storage"
	"url_shortener/pkg/logger"
)

func main() {
	// Инициализация логера
	logger.InitLogger()
	logger.Log.Info("Loading configuration...")
	config.LoadConfig()

	// Подключение к базе данных
	logger.Log.Info("Connecting to the storage...")
	storage.InitDB()
	defer storage.CloseDB()

	// Запуск сервера
	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", port)

	logger.Log.Info("Starting server on port ", port)

	// Определяем маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener is running"))
	})

	// Запуск сервера
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Log.Fatal("Server error: ", err)
	}
}
