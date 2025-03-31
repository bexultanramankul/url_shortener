package main

import (
	"url_shortener/pkg/logger"
)

func main() {
	// Инициализация логера
	logger.InitLogger()
	logger.Log.Info("Loading configuration...")
}
