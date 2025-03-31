package config

import (
	"strings"
	"url_shortener/pkg/logger"

	"github.com/spf13/viper"
)

// Config - основная структура конфигурации
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig - настройки сервера
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig - настройки базы данных
type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

var AppConfig Config

// LoadConfig загружает конфигурацию из config.yaml и переменных окружения
func LoadConfig() {
	log := logger.Log

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Настройка ENV-переменных
	viper.SetEnvPrefix("URL_SHORTENER") // Префикс для ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Читаем из файла, но не падаем, если его нет
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("config.yaml not found, loading from environment variables")
	}

	// Разбираем конфиг
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal("Failed to parse configuration: ", err)
	}

	// Проверяем обязательные поля
	if AppConfig.Database.User == "" || AppConfig.Database.Password == "" {
		log.Fatal("Missing required storage configuration (user/password)")
	}

	log.Info("Configuration loaded successfully")
}
