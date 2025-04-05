package main

import (
	"url_shortener/internal/config"
	"url_shortener/internal/server"
	"url_shortener/internal/storage"
	"url_shortener/pkg/logger"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("Loading configuration...")
	config.LoadConfig()

	logger.Log.Info("Connecting to the storage...")
	storage.InitDB()
	defer storage.CloseDB()

	logger.Log.Info("Connecting to the redis server...")
	storage.InitRedis()
	defer storage.CloseRedis()

	server := server.NewServer()
	server.Run()
}
